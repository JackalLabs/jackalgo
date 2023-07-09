package handler_wallet

import (
	"github.com/cosmos/cosmos-sdk/client"
)

type WalletHandler struct {
	clientCtx client.Context
	address   string
}
