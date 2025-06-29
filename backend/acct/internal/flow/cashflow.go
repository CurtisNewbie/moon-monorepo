package flow

import (
	"net/http"
	"os"

	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/money"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"gorm.io/gorm"
)

var (
	importPool    *util.AsyncPool = util.NewIOAsyncPool()
	categoryConfs map[string]CategoryConf
)

const (
	DirectionIn  = "IN"
	DirectionOut = "OUT"
)

type CategoryConf struct {
	Code string
	Name string
}

func LoadCategoryConfs(rail miso.Rail) {
	var cate []CategoryConf
	miso.UnmarshalFromPropKey("acct.category.builtin", &cate)
	categoryConfs = make(map[string]CategoryConf, len(cate))
	for i, v := range cate {
		categoryConfs[v.Code] = cate[i]
	}
	rail.Debugf("Loaded conf: %#v", categoryConfs)
}

type ListCashFlowReq struct {
	Paging         miso.Paging `desc:"Paging"`
	Direction      string      `desc:"Flow Direction: IN / OUT" valid:"member:IN|OUT|"`
	TransTimeStart *util.ETime `desc:"Transaction Time Range Start"`
	TransTimeEnd   *util.ETime `desc:"Transaction Time Range End"`
	TransId        string      `desc:"Transaction ID"`
	Category       string      `desc:"Category Code"`
	MinAmt         *money.Amt  `desc:"Minimum amount"`
}

type ListCashFlowRes struct {
	Direction     string     `desc:"Flow Direction: IN / OUT"`
	TransTime     util.ETime `desc:"Transaction Time"`
	TransId       string     `desc:"Transaction ID"`
	Counterparty  string     `desc:"Counterparty of the transaction"`
	PaymentMethod string     `desc:"Payment Method"`
	Amount        string     `desc:"Amount"`
	Currency      string     `desc:"Currency"`
	Extra         string     `desc:"Extra Information"`
	Category      string     `desc:"Category Code"`
	CategoryName  string     `desc:"Category Name"`
	Remark        string     `desc:"Remark"`
	CreatedAt     util.ETime `desc:"Create Time"`
}

func ListCashFlows(rail miso.Rail, db *gorm.DB, user common.User, req ListCashFlowReq) (miso.PageRes[ListCashFlowRes], error) {
	return dbquery.NewPagedQuery[ListCashFlowRes](db).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			q = q.Table(`cashflow`).
				Where("user_no = ?", user.UserNo).
				Where("deleted = 0")

			q = q.EqNotEmpty("trans_id", req.TransId).
				EqNotEmpty("category", req.Category).
				WhereNotNil("trans_time >= ?", req.TransTimeStart).
				WhereNotNil("trans_time <= ?", req.TransTimeEnd)

			if req.MinAmt != nil {
				abs := req.MinAmt.Abs()
				if abs.Cmp(money.Zero()) > 0 {
					q = q.Where("amount >= ?", abs)
					if req.MinAmt.Cmp(money.Zero()) < 0 {
						q = q.WhereIf(req.Direction != DirectionOut, "direction = ?", DirectionOut)
					} else {
						q = q.WhereIf(req.Direction != DirectionIn, "direction = ?", DirectionIn)
					}
				}
			}
			q = q.EqNotEmpty("direction = ?", req.Direction)
			return q
		}).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select("direction", "trans_time", "trans_id", "counterparty",
				"amount", "currency", "extra", "category", "remark", "created_at", "payment_method").
				Order("trans_time desc")
		}).
		Transform(func(t ListCashFlowRes) ListCashFlowRes {
			if v, ok := categoryConfs[t.Category]; ok {
				t.CategoryName = v.Name
			}
			t.Amount = money.UnitFmt(t.Amount, t.Currency)
			return t
		}).
		Scan(rail, req.Paging)
}

func ImportWechatCashflows(r *http.Request, rail miso.Rail, db *gorm.DB, user common.User) error {
	rail.Infof("User %v importing wechat cashflows", user.Username)

	path, err := util.SaveTmpFile("/tmp", r.Body)
	if err != nil {
		return err
	}
	rail.Infof("Wechat cashflows saved to temp file: %v", path)

	importPool.Go(func() {
		rail := rail.NextSpan()
		defer func() {
			os.Remove(path)
			rail.Infof("Temp file removed, %v", path)
		}()

		records, err := ParseWechatCashflows(rail, path)
		if err != nil {
			rail.Errorf("failed to parse wechat cashflows for %v, %v", user.Username, err)
			return
		}
		rail.Infof("Wechat cashflows (%d records) parsed for %v", len(records), user.Username)
		if len(records) > 0 {
			param := SaveCashflowParams{
				Cashflows: records,
				User:      user,
				Category:  WechatCategory,
			}
			saved, err := SaveCashflows(rail, db, param)
			if err != nil {
				rail.Errorf("failed to save wechat cashflows for %v, %v", user.Username, err)
			}

			changes := util.MapTo(saved, func(nc NewCashflow) CashflowChange { return CashflowChange{TransTime: nc.TransTime} })
			if err := OnCashflowChanged(rail, changes, user.UserNo); err != nil {
				rail.Errorf("Failed to update cashflow statistics for cashflow import, userNo: %v, %v", user.UserNo, err)
			}
		}
	})

	return nil
}

type NewCashflow struct {
	Direction     string
	TransTime     util.ETime
	TransId       string
	PaymentMethod string
	Counterparty  string
	Amount        string
	Currency      string
	Extra         string
	Remark        string
}

type SaveCashflowParams struct {
	Cashflows []NewCashflow
	Category  string
	User      common.User
}

type SavingCashflow struct {
	UserNo        string
	Direction     string
	TransTime     util.ETime
	TransId       string
	Counterparty  string
	Amount        string
	PaymentMethod string
	Currency      string
	Extra         string
	Category      string
	Remark        string
	CreatedAt     util.ETime
}

type CashflowCurrency struct {
	UserNo   string
	Currency string
}

func SaveCashflows(rail miso.Rail, db *gorm.DB, param SaveCashflowParams) ([]NewCashflow, error) {
	records := param.Cashflows
	if len(records) < 1 {
		return nil, nil
	}
	userNo := param.User.UserNo
	lock := userCashflowLock(rail, userNo)
	if err := lock.Lock(); err != nil {
		return nil, err
	}
	defer lock.Unlock()

	now := util.Now()

	// find those that already exist and skip them
	transIdSet := util.NewSet[string]()
	for _, v := range records {
		transIdSet.Add(v.TransId)
	}
	var existingTransId []string
	_, err := dbquery.NewQueryRail(rail, db).
		Raw(`SELECT trans_id FROM cashflow WHERE user_no = ? AND category = ? AND trans_id IN ? AND deleted = 0`,
			userNo, param.Category, transIdSet.CopyKeys()).
		Scan(&existingTransId)
	if err != nil {
		return nil, err
	}
	for _, ti := range existingTransId {
		rail.Debugf("Transaction %v (%v) for user %v already exists, ignored", ti, param.Category, userNo)
		transIdSet.Del(ti)
	}
	records = util.Filter(records, func(p NewCashflow) bool { return transIdSet.Has(p.TransId) })
	if len(records) < 1 {
		return nil, nil
	}

	ccySet := util.NewSet[string]()
	saving := make([]SavingCashflow, 0, len(records))
	for _, v := range records {
		transIdSet.Add(v.TransId)
		s := SavingCashflow{
			UserNo:        param.User.UserNo,
			Category:      param.Category,
			PaymentMethod: v.PaymentMethod,
			Direction:     v.Direction,
			TransTime:     v.TransTime,
			TransId:       v.TransId,
			Counterparty:  v.Counterparty,
			Amount:        v.Amount,
			Currency:      v.Currency,
			Extra:         v.Extra,
			Remark:        v.Remark,
			CreatedAt:     now,
		}
		saving = append(saving, s)
		ccySet.Add(v.Currency)
	}

	rail.Infof("Cashflows (%d records) saved for %v", len(saving), param.User.Username)
	_, err = dbquery.NewQueryRail(rail, db).Table("cashflow").Create(saving)
	if err != nil {
		return nil, err
	}

	newUserCcy := util.MapTo(ccySet.CopyKeys(), func(ccy string) CashflowCurrency { return CashflowCurrency{UserNo: userNo, Currency: ccy} })
	_, err = dbquery.NewQueryRail(rail, db).Table("cashflow_currency").CreateIgnore(newUserCcy)
	return records, err
}

func userCashflowLock(rail miso.Rail, userNo string) *redis.RLock {
	return redis.NewRLockf(rail, "acct:cashflow:user:%v", userNo)
}

func ListCurrencies(rail miso.Rail, db *gorm.DB, user common.User) ([]string, error) {
	var ccy []string
	_, err := dbquery.NewQueryRail(rail, db).
		Raw(`SELECT currency FROM cashflow_currency WHERE user_no = ?`, user.UserNo).
		Scan(&ccy)
	return ccy, err
}
