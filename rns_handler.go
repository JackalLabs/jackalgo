package jackalgo

import (
	"fmt"
)

func NewRnsHandler(w WalletHandler) *RnsHandler {

	r := RnsHandler{
		walletHandler: w,
	}

	return &r
}

func (r *RnsHandler) SayHello() {
	fmt.Println("Hello from RnsHandler")
}
