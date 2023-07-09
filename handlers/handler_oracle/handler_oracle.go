package handler_oracle

import (
	"fmt"

	"github.com/JackalLabs/jackalgo/handlers/handler_wallet"
)

func NewOracleHandler(w *handler_wallet.WalletHandler) *OracleHandler {

	o := OracleHandler{
		walletHandler: w,
	}

	return &o
}

func (o *OracleHandler) SayHello() {
	fmt.Println("Hello from OracleHandler")
}
