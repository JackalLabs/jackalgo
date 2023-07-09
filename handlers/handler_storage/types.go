package handler_storage

import (
	"github.com/JackalLabs/jackalgo/handlers/handler_wallet"
)

type StorageHandler struct {
	walletHandler *handler_wallet.WalletHandler
}

func NewStorageHandler(w *handler_wallet.WalletHandler) *StorageHandler {

	s := StorageHandler{
		walletHandler: w,
	}

	return &s
}
