package gov_handler

import "github.com/JackalLabs/jackalgo/handlers/wallet_handler"

type GovHandler struct {
	walletHandler *wallet_handler.WalletHandler
}

func NewGovHandler(w *wallet_handler.WalletHandler) *GovHandler {
	g := GovHandler{
		walletHandler: w,
	}

	return &g
}
