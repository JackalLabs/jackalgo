package handler_oracle

import (
	"fmt"
	"github.com/JackalLabs/jackalgo"
)

func NewOracleHandler(w *jackalgo.WalletHandler) *OracleHandler {

	o := OracleHandler{
		walletHandler: w,
	}

	return &o
}

func (o *OracleHandler) SayHello() {
	fmt.Println("Hello from OracleHandler")
}
