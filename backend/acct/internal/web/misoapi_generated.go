// auto generated by misoapi v0.1.18 at 2025/04/06 22:18:31 (CST), please do not modify
package web

import (
	"github.com/curtisnewbie/acct/internal/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
)

func init() {
	miso.IPost("/open/api/v1/cashflow/list",
		func(inb *miso.Inbound, req flow.ListCashFlowReq) (miso.PageRes[flow.ListCashFlowRes], error) {
			return ApiListCashFlows(inb.Rail(), dbquery.GetDB(), common.GetUser(inb.Rail()), req)
		}).
		Extra(miso.ExtraName, "ApiListCashFlows").
		Resource(CodeManageCashflows)

	miso.Post("/open/api/v1/cashflow/import/wechat",
		func(inb *miso.Inbound) (any, error) {
			return nil, ApiImportWechatCashflows(inb, inb.Rail(), dbquery.GetDB(), common.GetUser(inb.Rail()))
		}).
		Extra(miso.ExtraName, "ApiImportWechatCashflows").
		Resource(CodeManageCashflows)

	miso.Get("/open/api/v1/cashflow/list-currency",
		func(inb *miso.Inbound) ([]string, error) {
			return ApiListCurrency(inb.Rail(), dbquery.GetDB(), common.GetUser(inb.Rail()))
		}).
		Extra(miso.ExtraName, "ApiListCurrency").
		Resource(CodeManageCashflows)

	miso.IPost("/open/api/v1/cashflow/list-statistics",
		func(inb *miso.Inbound, req flow.ApiListStatisticsReq) (miso.PageRes[flow.ApiListStatisticsRes], error) {
			return ApiListCashflowStatistics(inb.Rail(), dbquery.GetDB(), common.GetUser(inb.Rail()), req)
		}).
		Extra(miso.ExtraName, "ApiListCashflowStatistics").
		Resource(CodeManageCashflows)

	miso.IPost("/open/api/v1/cashflow/plot-statistics",
		func(inb *miso.Inbound, req flow.ApiPlotStatisticsReq) ([]flow.ApiPlotStatisticsRes, error) {
			return ApiPlotCashflowStatistics(inb.Rail(), dbquery.GetDB(), common.GetUser(inb.Rail()), req)
		}).
		Extra(miso.ExtraName, "ApiPlotCashflowStatistics").
		Resource(CodeManageCashflows)

}
