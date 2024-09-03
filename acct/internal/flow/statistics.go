package flow

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/curtisnewbie/miso/middleware/money"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"gorm.io/gorm"
)

const (
	AggTypeYearly  = "YEARLY"
	AggTypeMonthly = "MONTHLY"
	AggTypeWeekly  = "WEEKLY"
)

var (
	RangeFormatMap = map[string]string{
		AggTypeYearly:  `2006`,
		AggTypeMonthly: `200601`,
		AggTypeWeekly:  `20060102`,
	}

	CalcAggStatPipeline = rabbit.NewEventPipeline[CalcCashflowStatsEvent]("acct:cashflow:calc-agg-stat").
				LogPayload().
				Listen(2, OnCalcCashflowStatsEvent).
				MaxRetry(3)
)

type CalcCashflowStatsEvent struct {
	UserNo   string
	AggType  string
	AggRange string
	AggTime  util.ETime
}

type ApiCalcCashflowStatsReq struct {
	AggType  string `desc:"Aggregation Type." valid:"member:YEARLY|MONTHLY|WEEKLY"`
	AggRange string `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)." valid:"notEmpty"`
}

func ParseAggRangeTime(aggType string, aggRange string) (util.ETime, error) {
	pat, ok := RangeFormatMap[aggType]
	if !ok {
		return util.ETime{}, miso.NewErrf("Invalid AggType")
	}

	t, err := time.ParseInLocation(pat, aggRange, time.Local)
	if err != nil {
		return util.ETime{}, miso.NewErrf("Invalid AppRange '%s' for %s aggregate type", aggRange, aggType).
			WithInternalMsg("%v", err)
	}
	if aggType == AggTypeWeekly {
		wd := t.Weekday()
		if wd != time.Sunday {
			return util.ETime{}, miso.NewErrf("Invalid aggRange '%v' for aggType: %v, should be Sunday", aggRange, aggType)
		}
	}
	return util.ToETime(t), err
}

type CashflowChange struct {
	TransTime util.ETime
}

func OnCashflowChanged(rail miso.Rail, changes []CashflowChange, userNo string) error {
	if len(changes) < 1 {
		return nil
	}

	aggMap := map[string]util.Set[string]{}
	mapAddAgg := func(typ, val string) {
		prev, ok := aggMap[typ]
		if !ok {
			v := util.NewSet[string]()
			aggMap[typ] = v
			prev = v
		}
		prev.Add(val)
	}

	for _, c := range changes {
		tt := c.TransTime.ToTime()
		mapAddAgg(AggTypeYearly, tt.Format(RangeFormatMap[AggTypeYearly]))
		mapAddAgg(AggTypeMonthly, tt.Format(RangeFormatMap[AggTypeMonthly]))
		mapAddAgg(AggTypeWeekly, tt.AddDate(0, 0, -(int(tt.Weekday())-int(time.Sunday))).Format(RangeFormatMap[AggTypeWeekly]))
	}

	for typ, set := range aggMap {
		for val := range set.Keys {
			err := CalcCashflowStatsAsync(rail, ApiCalcCashflowStatsReq{AggType: typ, AggRange: val}, userNo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CalcCashflowStatsAsync(rail miso.Rail, req ApiCalcCashflowStatsReq, userNo string) error {
	t, err := ParseAggRangeTime(req.AggType, req.AggRange)
	if err != nil {
		return err
	}
	rail = rail.NextSpan()
	return CalcAggStatPipeline.Send(rail, CalcCashflowStatsEvent{
		AggType:  req.AggType,
		AggRange: req.AggRange,
		AggTime:  t,
		UserNo:   userNo,
	})
}

func OnCalcCashflowStatsEvent(rail miso.Rail, evt CalcCashflowStatsEvent) error {
	rlock := miso.NewRLockf(rail, "acct:calc-cashflow-stats:%v:%v:%v", evt.UserNo, evt.AggType, evt.AggRange)
	if err := rlock.Lock(); err != nil {
		return err
	}
	defer rlock.Unlock()

	db := miso.GetMySQL()
	t := evt.AggTime.ToTime()
	switch evt.AggType {
	case AggTypeMonthly:
		return calcMonthlyCashflow(rail, db, t, evt.AggRange, evt.UserNo)
	case AggTypeWeekly:
		return calcWeeklyCashflow(rail, db, t, evt.AggRange, evt.UserNo)
	case AggTypeYearly:
		return calcYearlyCashflow(rail, db, t, evt.AggRange, evt.UserNo)
	}
	return nil
}

func calcYearlyCashflow(rail miso.Rail, db *gorm.DB, t time.Time, aggRange string, userNo string) error {
	start := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	lastDay := time.Date(t.Year(), 12, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1)
	end := time.Date(t.Year(), 12, lastDay.Day(), 23, 59, 59, 0, time.Local)
	sum, err := calcCashflowSum(rail, db, TimeRange{Start: start, End: end}, userNo)
	if err != nil {
		return err
	}
	return updateCashflowStat(rail, db, sum, AggTypeYearly, aggRange, userNo)
}

func calcMonthlyCashflow(rail miso.Rail, db *gorm.DB, t time.Time, aggRange string, userNo string) error {
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	lastDay := time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1)
	end := time.Date(t.Year(), t.Month(), lastDay.Day(), 23, 59, 59, 0, time.Local)
	sum, err := calcCashflowSum(rail, db, TimeRange{Start: start, End: end}, userNo)
	if err != nil {
		return err
	}
	return updateCashflowStat(rail, db, sum, AggTypeMonthly, aggRange, userNo)
}

func calcWeeklyCashflow(rail miso.Rail, db *gorm.DB, t time.Time, aggRange string, userNo string) error {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local) // sunday
	lastDay := t.AddDate(0, 0, 6)
	end := time.Date(t.Year(), lastDay.Month(), lastDay.Day(), 23, 59, 59, 0, time.Local)
	sum, err := calcCashflowSum(rail, db, TimeRange{Start: start, End: end}, userNo)
	if err != nil {
		return err
	}
	return updateCashflowStat(rail, db, sum, AggTypeWeekly, aggRange, userNo)
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type CashflowSum struct {
	Currency  string
	AmountSum string
}

func calcCashflowSum(rail miso.Rail, db *gorm.DB, tr TimeRange, userNo string) ([]CashflowSum, error) {
	if tr.Start.After(tr.End) {
		tr.Start, tr.End = tr.End, tr.Start
	}
	rail.Infof("Calculating cashflow sum between %v, %v, userNo: %v", tr.Start, tr.End, userNo)

	var res []CashflowSum
	err := db.Raw(`
	SELECT SUM(case when direction = 'IN' then amount else (-1 * amount) end) amount_sum, currency
	FROM cashflow WHERE user_no = ? and trans_time between ? and ? and deleted = 0
	GROUP BY currency
	`,
		userNo, tr.Start, tr.End).
		Scan(&res).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to query cashflow sum, %w", err)
	}
	return res, nil
}

func updateCashflowStat(rail miso.Rail, db *gorm.DB, stats []CashflowSum, aggType string, aggRange string, userNo string) error {
	for _, st := range stats {
		var id int64
		err := db.Raw(`SELECT id FROM cashflow_statistics WHERE user_no = ? and agg_type = ? and agg_range = ? and currency = ?`,
			userNo, aggType, aggRange, st.Currency).Scan(&id).Error
		if err != nil {
			return fmt.Errorf("failed to query cashflow_statistics, %w", err)
		}
		if id > 0 {
			err := db.Exec(`UPDATE cashflow_statistics SET agg_value = ? WHERE id = ?`,
				st.AmountSum, id).Error
			if err != nil {
				return fmt.Errorf("failed to update cashflow_statistics, id: %v, %w", id, err)
			}
		} else {
			err := db.Exec(`INSERT INTO cashflow_statistics (user_no, agg_type, agg_range, currency, agg_value) VALUES (?,?,?,?,?)`,
				userNo, aggType, aggRange, st.Currency, st.AmountSum).Error
			if err != nil {
				return fmt.Errorf("failed to save cashflow_statistics, %w", err)
			}
		}
	}
	return nil
}

type ApiListStatisticsReq struct {
	Paging   miso.Paging `desc:"Paging Info"`
	AggType  string      `desc:"Aggregation Type." valid:"member:YEARLY|MONTHLY|WEEKLY"`
	AggRange string      `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)."`
	Currency string      `desc:"Currency"`
}

type ApiListStatisticsRes struct {
	AggType  string `desc:"Aggregation Type."`
	AggRange string `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)."`
	AggValue string `desc:"Aggregation Value."`
	Currency string `desc:"Currency"`
}

func ListCashflowStatistics(rail miso.Rail, db *gorm.DB, req ApiListStatisticsReq, user common.User) (miso.PageRes[ApiListStatisticsRes], error) {

	if req.AggRange != "" {
		_, err := ParseAggRangeTime(req.AggType, req.AggRange)
		if err != nil {
			return miso.PageRes[ApiListStatisticsRes]{}, err
		}
	}

	return miso.NewPageQuery[ApiListStatisticsRes]().
		WithPage(req.Paging).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Table(`cashflow_statistics`).
				Where(`user_no = ?`, user.UserNo).
				Where(`agg_type = ?`, req.AggType).
				Order("agg_range desc, currency desc")
			if req.AggRange != "" {
				tx = tx.Where("agg_range = ?", req.AggRange)
			}
			if req.Currency != "" {
				tx = tx.Where("currency = ?", req.Currency)
			}
			return tx
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select("agg_type, agg_range, agg_value, currency")
		}).
		ForEach(func(t ApiListStatisticsRes) ApiListStatisticsRes {
			t.AggValue = money.UnitFmt(t.AggValue, t.Currency)
			return t
		}).
		Exec(rail, db)
}

type ApiPlotStatisticsReq struct {
	StartTime util.ETime `desc:"Start time"`
	EndTime   util.ETime `desc:"End time"`
	AggType   string     `desc:"Aggregation Type." valid:"member:YEARLY|MONTHLY|WEEKLY"`
	Currency  string     `desc:"Currency"`
}

type ApiPlotStatisticsRes struct {
	AggRange string `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)."`
	AggValue string `desc:"Aggregation Value."`
}

func PlotCashflowStatistics(rail miso.Rail, db *gorm.DB, req ApiPlotStatisticsReq, user common.User) ([]ApiPlotStatisticsRes, error) {
	if req.StartTime.After(req.EndTime) {
		req.StartTime, req.EndTime = req.EndTime, req.StartTime
	}

	var pad string = ""
	var res []ApiPlotStatisticsRes
	switch req.AggType {
	case AggTypeMonthly:
		pad = "01"
	case AggTypeYearly:
		pad = "0101"
	}

	err := db.Raw(`
			SELECT agg_range, agg_value FROM cashflow_statistics
			WHERE user_no = ? AND agg_type = ? AND currency = ?
			AND str_to_date(concat(agg_range, ?), '%Y%m%d') BETWEEN ? AND ?`,
		user.UserNo, req.AggType, req.Currency, pad, req.StartTime, req.EndTime).Scan(&res).Error
	if err == nil {
		if res == nil {
			res = []ApiPlotStatisticsRes{}
		}
		set := util.NewSet[string]()
		for _, r := range res {
			set.Add(r.AggRange)
		}
		start := req.StartTime
		for start.Before(req.EndTime) {
			var next string
			switch req.AggType {
			case AggTypeYearly:
				next = start.Format(RangeFormatMap[AggTypeYearly])
			case AggTypeMonthly:
				next = start.Format(RangeFormatMap[AggTypeMonthly])
			case AggTypeWeekly:
				sun := start.AddDate(0, 0, -(int(start.Weekday()) - int(time.Sunday)))
				if !sun.Before(start) {
					next = start.Format(RangeFormatMap[AggTypeWeekly])
				} else {
					start = sun
				}
			}

			if next != "" && set.Add(next) {
				res = append(res, ApiPlotStatisticsRes{AggRange: next, AggValue: "0"})
			}

			switch req.AggType {
			case AggTypeYearly:
				start = start.AddDate(1, 0, 0)
			case AggTypeMonthly:
				start = start.AddDate(0, 1, 0)
			case AggTypeWeekly:
				start = start.AddDate(0, 0, 7)
			}
		}
		sort.Slice(res, func(i, j int) bool { return strings.Compare(res[i].AggRange, res[j].AggRange) < 0 })
	}
	return res, err
}
