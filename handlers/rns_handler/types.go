package rns_handler

import "github.com/JackalLabs/jackalgo/handlers/wallet_handler"

type RnsHandler struct {
	walletHandler *wallet_handler.WalletHandler
}

func NewRnsHandler(w *wallet_handler.WalletHandler) *RnsHandler {
	r := RnsHandler{
		walletHandler: w,
	}

	return &r
}
