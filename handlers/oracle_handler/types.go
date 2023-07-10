package oracle_handler

import "github.com/JackalLabs/jackalgo/handlers/wallet_handler"

type OracleHandler struct {
	walletHandler *wallet_handler.WalletHandler
}

func NewOracleHandler(w *wallet_handler.WalletHandler) *OracleHandler {

	o := OracleHandler{
		walletHandler: w,
	}

	return &o
}
