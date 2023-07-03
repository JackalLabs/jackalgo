package jackalgo

import (
	"github.com/cosmos/cosmos-sdk/client"
)

type WalletHandler struct {
	clientCtx client.Context
	address   string
}

type RnsHandler struct {
	walletHandler *WalletHandler
}
