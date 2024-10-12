package web

import (
	"github.com/curtisnewbie/acct/internal/flow"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

const (
	CodeManageCashflows = "acct:ManageCashflows"
)

func RegisterEndpoints(rail miso.Rail) {
	common.LoadBuiltinPropagationKeys()
	auth.ExposeResourceInfo([]auth.Resource{{
		Code: CodeManageCashflows,
		Name: "Manage Personal Cashflows",
	}})
}

// misoapi-http: POST /open/api/v1/cashflow/list
// misoapi-resource: ref(CodeManageCashflows)
func ApiListCashFlows(rail miso.Rail, db *gorm.DB, user common.User, req flow.ListCashFlowReq) (miso.PageRes[flow.ListCashFlowRes], error) {
	return flow.ListCashFlows(rail, db, user, req)
}

// misoapi-http: POST /open/api/v1/cashflow/import/wechat
// misoapi-resource: ref(CodeManageCashflows)
func ApiImportWechatCashflows(inb *miso.Inbound, rail miso.Rail, db *gorm.DB, user common.User) error {
	_, r := inb.Unwrap()
	return flow.ImportWechatCashflows(r, rail, db, user)
}

// misoapi-http: GET /open/api/v1/cashflow/list-currency
// misoapi-resource: ref(CodeManageCashflows)
func ApiListCurrency(rail miso.Rail, db *gorm.DB, user common.User) ([]string, error) {
	return flow.ListCurrencies(rail, mysql.GetMySQL(), user)
}

// misoapi-http: POST /open/api/v1/cashflow/list-statistics
// misoapi-resource: ref(CodeManageCashflows)
func ApiListCashflowStatistics(rail miso.Rail, db *gorm.DB, user common.User, req flow.ApiListStatisticsReq) (miso.PageRes[flow.ApiListStatisticsRes], error) {
	return flow.ListCashflowStatistics(rail, mysql.GetMySQL(), req, user)
}

// misoapi-http: POST /open/api/v1/cashflow/plot-statistics
// misoapi-resource: ref(CodeManageCashflows)
func ApiPlotCashflowStatistics(rail miso.Rail, db *gorm.DB, user common.User, req flow.ApiPlotStatisticsReq) ([]flow.ApiPlotStatisticsRes, error) {
	return flow.PlotCashflowStatistics(rail, db, req, user)
}
