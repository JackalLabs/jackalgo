package jackalgo

import (
	"context"
	"fmt"
	"github.com/JackalLabs/jackalgo/handler_file_io"
	"github.com/JackalLabs/jackalgo/handler_gov"
	"github.com/JackalLabs/jackalgo/handler_oracle"
	"github.com/JackalLabs/jackalgo/handler_rns"
	"github.com/JackalLabs/jackalgo/handler_storage"

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

func (w *WalletHandler) NewFileIoHandler() *handler_file_io.FileIoHandler {
	return handler_file_io.NewFileIoHandler(w)
}
func (w *WalletHandler) NewGovHandler() *handler_gov.GovHandler {
	return handler_gov.NewGovHandler(w)
}
func (w *WalletHandler) NewOracleHandler() *handler_oracle.OracleHandler {
	return handler_oracle.NewOracleHandler(w)
}
func (w *WalletHandler) NewRnsHandler() *handler_rns.RnsHandler {
	return handler_rns.NewRnsHandler(w)
}
func (w *WalletHandler) NewStorageHandler() *handler_storage.StorageHandler {
	return handler_storage.NewStorageHandler(w)
}

func (w *WalletHandler) GetChainID() string {
	return w.clientCtx.ChainID
}

func (w *WalletHandler) SendTokens(toAddress string, amount types.Coins) (*types.TxResponse, error) {

	sendMsg := banktypes.MsgSend{
		FromAddress: w.address,
		ToAddress:   toAddress,
		Amount:      amount,
	}

	res, err := tx.SendTx(w.clientCtx, nil, &sendMsg)

	return res, err
}
