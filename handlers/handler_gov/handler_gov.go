package handler_gov

import (
	"fmt"

	"github.com/JackalLabs/jackalgo/handlers/handler_wallet"
)

func NewGovHandler(w *handler_wallet.WalletHandler) *GovHandler {

	g := GovHandler{
		walletHandler: w,
	}

	return &g
}

func (g *GovHandler) SayHello() {
	fmt.Println("Hello from GovHandler")
}
