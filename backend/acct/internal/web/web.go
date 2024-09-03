package web

import (
	"github.com/curtisnewbie/acct/internal/flow"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
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

	miso.GroupRoute("/open/api/v1",
		miso.IPost("/cashflow/list", ApiListCashFlows).Resource(CodeManageCashflows),
		miso.Post("/cashflow/import/wechat", ApiImportWechatCashflows).Resource(CodeManageCashflows),
		miso.Get("/cashflow/list-currency", ApiListCurrency).Resource(CodeManageCashflows),
		miso.IPost("/cashflow/list-statistics", ApiListCashflowStatistics).Resource(CodeManageCashflows),
		miso.IPost("/cashflow/plot-statistics", ApiPlotCashflowStatistics).Resource(CodeManageCashflows),
	)
}

func ApiListCashFlows(inb *miso.Inbound, req flow.ListCashFlowReq) (miso.PageRes[flow.ListCashFlowRes], error) {
	return flow.ListCashFlows(inb.Rail(), miso.GetMySQL(), common.GetUser(inb.Rail()), req)
}

func ApiImportWechatCashflows(inb *miso.Inbound) (any, error) {
	return nil, flow.ImportWechatCashflows(inb, miso.GetMySQL())
}

func ApiListCashflowStatistics(inb *miso.Inbound, req flow.ApiListStatisticsReq) (miso.PageRes[flow.ApiListStatisticsRes], error) {
	return flow.ListCashflowStatistics(inb.Rail(), miso.GetMySQL(), req, common.GetUser(inb.Rail()))
}

func ApiListCurrency(inb *miso.Inbound) ([]string, error) {
	return flow.ListCurrencies(inb.Rail(), miso.GetMySQL(), common.GetUser(inb.Rail()))
}

func ApiPlotCashflowStatistics(inb *miso.Inbound, req flow.ApiPlotStatisticsReq) ([]flow.ApiPlotStatisticsRes, error) {
	return flow.PlotCashflowStatistics(inb.Rail(), miso.GetMySQL(), req, common.GetUser(inb.Rail()))
}
