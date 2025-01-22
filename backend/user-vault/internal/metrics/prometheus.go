package metrics

import "github.com/curtisnewbie/miso/miso"

var (
	FetchUserInfoHisto       = miso.NewPromHisto("user_vault_fetch_user_info_duration")
	TokenExchangeHisto       = miso.NewPromHisto("user_vault_token_exchange_duration")
	ResourceAccessCheckHisto = miso.NewPromHisto("user_vault_resource_access_check_duration")
)
