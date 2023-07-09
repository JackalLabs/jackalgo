package handler_rns

import (
	"fmt"

	"github.com/JackalLabs/jackalgo/handlers/handler_wallet"
)

func NewRnsHandler(w *handler_wallet.WalletHandler) *RnsHandler {

	r := RnsHandler{
		walletHandler: w,
	}

	return &r
}

func (r *RnsHandler) SayHello() {
	fmt.Println("Hello from RnsHandler")
}
