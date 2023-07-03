package jackalgo

import (
	"context"

	"github.com/JackalLabs/jackalgo/utils"
	"github.com/cosmos/cosmos-sdk/client"
)

func NewWalletHandler() *WalletHandler {

	k := &client.Context{}

	srvCtx := utils.NewDefaultContext()
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, k)
	ctx = context.WithValue(ctx, utils.JackalGoContextKey, srvCtx)

	var clientCtx client.Context

	if v := ctx.Value(client.ClientContextKey); v != nil {
		clientCtxPtr := v.(*client.Context)
		clientCtx = *clientCtxPtr
	}

	clientCtx = clientCtx.WithChainID("jackal-1")

	w := WalletHandler{
		clientCtx: clientCtx,
	}

	return &w
}

func (w *WalletHandler) GetChainID() string {
	return w.clientCtx.ChainID
}
