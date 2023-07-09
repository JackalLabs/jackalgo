package handler_wallet

import (
	"context"
	"fmt"

	"github.com/JackalLabs/jackalgo/tx"
	"github.com/JackalLabs/jackalgo/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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

func (w *WalletHandler) SayHello() {
	fmt.Println("Hello from WalletHandler")
}

func (w *WalletHandler) GetChainID() string {
	return w.clientCtx.ChainID
}

func (w *WalletHandler) GetAddress() string {
	return w.address
}

func (w *WalletHandler) GetClientCtx() client.Context {
	return w.clientCtx
}

func (w *WalletHandler) SendTx(msg types.Msg) (*types.TxResponse, error) {

	res, err := tx.SendTx(w.clientCtx, nil, msg)

	return res, err
}

func (w *WalletHandler) SendTokens(toAddress string, amount types.Coins) (*types.TxResponse, error) {

	sendMsg := banktypes.MsgSend{
		FromAddress: w.address,
		ToAddress:   toAddress,
		Amount:      amount,
	}

	res, err := w.SendTx(&sendMsg)

	return res, err
}
