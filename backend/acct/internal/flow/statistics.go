package flow

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/money"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/errs"
	"github.com/curtisnewbie/miso/util/hash"
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
	AggTime  atom.Time
}

type ApiCalcCashflowStatsReq struct {
	AggType  string `desc:"Aggregation Type." valid:"member:YEARLY|MONTHLY|WEEKLY"`
	AggRange string `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)." valid:"notEmpty"`
}

func ParseAggRangeTime(aggType string, aggRange string) (atom.Time, error) {
	pat, ok := RangeFormatMap[aggType]
	if !ok {
		return atom.Time{}, errs.NewErrf("Invalid AggType")
	}

	t, err := time.ParseInLocation(pat, aggRange, time.Local)
	if err != nil {
		return atom.Time{}, errs.NewErrf("Invalid AppRange '%s' for %s aggregate type", aggRange, aggType).
			WithInternalMsg("%v", err)
	}
	if aggType == AggTypeWeekly {
		wd := t.Weekday()
		if wd != time.Sunday {
			return atom.Time{}, errs.NewErrf("Invalid aggRange '%v' for aggType: %v, should be Sunday", aggRange, aggType)
		}
	}
	return atom.WrapTime(t), err
}

type CashflowChange struct {
	TransTime atom.Time
}

func OnCashflowChanged(rail miso.Rail, changes []CashflowChange, userNo string) error {
	if len(changes) < 1 {
		return nil
	}

	aggMap := map[string]hash.Set[string]{}
	mapAddAgg := func(typ, val string) {
		prev, ok := aggMap[typ]
		if !ok {
			v := hash.NewSet[string]()
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
	rlock := redis.NewRLockf(rail, "acct:calc-cashflow-stats:%v:%v:%v", evt.UserNo, evt.AggType, evt.AggRange)
	if err := rlock.Lock(); err != nil {
		return err
	}
	defer rlock.Unlock()

	db := mysql.GetMySQL()
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
	_, err := dbquery.NewQuery(rail, db).Raw(`
	SELECT SUM(case when direction = 'IN' then amount else (-1 * amount) end) amount_sum, currency
	FROM cashflow WHERE user_no = ? and trans_time between ? and ? and deleted = 0
	GROUP BY currency
	`, userNo, tr.Start, tr.End).Scan(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to query cashflow sum, %w", err)
	}
	return res, nil
}

func updateCashflowStat(rail miso.Rail, db *gorm.DB, stats []CashflowSum, aggType string, aggRange string, userNo string) error {
	for _, st := range stats {
		var id int64
		_, err := dbquery.NewQuery(rail, db).Raw(`SELECT id FROM cashflow_statistics WHERE user_no = ? and agg_type = ? and agg_range = ? and currency = ?`,
			userNo, aggType, aggRange, st.Currency).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to query cashflow_statistics, %w", err)
		}
		if id > 0 {
			_, err := dbquery.NewQuery(rail, db).Exec(`UPDATE cashflow_statistics SET agg_value = ? WHERE id = ?`,
				st.AmountSum, id)
			if err != nil {
				return fmt.Errorf("failed to update cashflow_statistics, id: %v, %w", id, err)
			}
		} else {
			_, err := dbquery.NewQuery(rail, db).Exec(`INSERT INTO cashflow_statistics (user_no, agg_type, agg_range, currency, agg_value) VALUES (?,?,?,?,?)`,
				userNo, aggType, aggRange, st.Currency, st.AmountSum)
			if err != nil {
				return fmt.Errorf("failed to save cashflow_statistics, %w", err)
			}
		}
	}
	return nil
}

type ApiListStatisticsReq struct {
	Paging   miso.Paging `desc:"Paging Info" json:"paging"`
	AggType  string      `desc:"Aggregation Type." valid:"member:YEARLY|MONTHLY|WEEKLY" json:"aggType"`
	AggRange string      `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)." json:"aggRange"`
	Currency string      `desc:"Currency" json:"currency"`
}

type ApiListStatisticsRes struct {
	AggType  string `desc:"Aggregation Type." json:"aggType"`
	AggRange string `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)." json:"aggRange"`
	AggValue string `desc:"Aggregation Value." json:"aggValue"`
	Currency string `desc:"Currency" json:"currency"`
}

func ListCashflowStatistics(rail miso.Rail, db *gorm.DB, req ApiListStatisticsReq, user common.User) (miso.PageRes[ApiListStatisticsRes], error) {

	if req.AggRange != "" {
		_, err := ParseAggRangeTime(req.AggType, req.AggRange)
		if err != nil {
			return miso.PageRes[ApiListStatisticsRes]{}, err
		}
	}

	return dbquery.NewPagedQuery[ApiListStatisticsRes](db).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			q = q.Table(`cashflow_statistics`).
				Eq(`user_no`, user.UserNo).
				Eq(`agg_type`, req.AggType).
				EqNotEmpty("agg_range", req.AggRange).
				EqNotEmpty("currency", req.Currency).
				Order("agg_range desc, currency desc")
			return q
		}).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select("agg_type, agg_range, agg_value, currency")
		}).
		Transform(func(t ApiListStatisticsRes) ApiListStatisticsRes {
			t.AggValue = money.UnitFmt(t.AggValue, t.Currency)
			return t
		}).
		Scan(rail, req.Paging)
}

type ApiPlotStatisticsReq struct {
	StartTime atom.Time `desc:"Start time" json:"startTime"`
	EndTime   atom.Time `desc:"End time" json:"endTime"`
	AggType   string    `desc:"Aggregation Type." valid:"member:YEARLY|MONTHLY|WEEKLY" json:"aggType"`
	Currency  string    `desc:"Currency" json:"currency"`
}

type ApiPlotStatisticsRes struct {
	AggRange string `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)." json:"aggRange"`
	AggValue string `desc:"Aggregation Value." json:"aggValue"`
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

	_, err := dbquery.NewQuery(rail, db).
		Raw(`
			SELECT agg_range, agg_value FROM cashflow_statistics
			WHERE user_no = ? AND agg_type = ? AND currency = ?
			AND str_to_date(concat(agg_range, ?), '%Y%m%d') BETWEEN ? AND ?`,
			user.UserNo, req.AggType, req.Currency, pad, req.StartTime, req.EndTime).Scan(&res)
	if err == nil {
		if res == nil {
			res = []ApiPlotStatisticsRes{}
		}
		set := hash.NewSet[string]()
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
