package handler_gov

import (
	"fmt"
	"github.com/JackalLabs/jackalgo"
)

func NewGovHandler(w *jackalgo.WalletHandler) *GovHandler {

	g := GovHandler{
		walletHandler: w,
	}

	return &g
}

func (g *GovHandler) SayHello() {
	fmt.Println("Hello from GovHandler")
}
