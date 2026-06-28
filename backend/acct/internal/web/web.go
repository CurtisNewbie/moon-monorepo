package web

import (
	"github.com/curtisnewbie/acct/internal/mflow"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

const (
	CodeManageCashflows = "acct:ManageCashflows"
)

func PrepareWebServer(rail miso.Rail) {
	RegisterApi()

	auth.ExposeResourceInfo([]auth.Resource{{
		Code: CodeManageCashflows,
		Name: "Manage Personal Cashflows",
	}})
}

// misoapi-http: POST /open/api/v1/cashflow/list
// misoapi-resource: ref(CodeManageCashflows)
func ApiListCashFlows(rail miso.Rail, db *gorm.DB, user flow.User, req mflow.ListCashFlowReq) (miso.PageRes[mflow.ListCashFlowRes], error) {
	return mflow.ListCashFlows(rail, db, user, req)
}

// misoapi-http: POST /open/api/v1/cashflow/import/wechat
// misoapi-resource: ref(CodeManageCashflows)
func ApiImportWechatCashflows(inb *miso.Inbound, rail miso.Rail, db *gorm.DB, user flow.User) error {
	_, r := inb.Unwrap()
	return mflow.ImportWechatCashflows(r, rail, db, user)
}

// misoapi-http: GET /open/api/v1/cashflow/list-currency
// misoapi-resource: ref(CodeManageCashflows)
func ApiListCurrency(rail miso.Rail, db *gorm.DB, user flow.User) ([]string, error) {
	return mflow.ListCurrencies(rail, mysql.GetMySQL(), user)
}

// misoapi-http: POST /open/api/v1/cashflow/list-statistics
// misoapi-resource: ref(CodeManageCashflows)
func ApiListCashflowStatistics(rail miso.Rail, db *gorm.DB, user flow.User, req mflow.ApiListStatisticsReq) (miso.PageRes[mflow.ApiListStatisticsRes], error) {
	return mflow.ListCashflowStatistics(rail, mysql.GetMySQL(), req, user)
}

// misoapi-http: POST /open/api/v1/cashflow/plot-statistics
// misoapi-resource: ref(CodeManageCashflows)
func ApiPlotCashflowStatistics(rail miso.Rail, db *gorm.DB, user flow.User, req mflow.ApiPlotStatisticsReq) ([]mflow.ApiPlotStatisticsRes, error) {
	return mflow.PlotCashflowStatistics(rail, db, req, user)
}
