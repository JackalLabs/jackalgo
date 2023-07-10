package wallet_handler

import (
	"github.com/cosmos/cosmos-sdk/client"
)

func (w *WalletHandler) GetChainID() string {
	return w.clientCtx.ChainID
}

func (w *WalletHandler) GetAddress() string {
	return w.address
}

func (w *WalletHandler) GetClientCtx() client.Context {
	return w.clientCtx
}
