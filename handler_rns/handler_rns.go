package handler_rns

import (
	"fmt"
	"github.com/JackalLabs/jackalgo"
)

func NewRnsHandler(w *jackalgo.WalletHandler) *RnsHandler {

	r := RnsHandler{
		walletHandler: w,
	}

	return &r
}

func (r *RnsHandler) SayHello() {
	fmt.Println("Hello from RnsHandler")
}
