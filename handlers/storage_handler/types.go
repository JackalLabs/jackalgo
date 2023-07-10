package storage_handler

import "github.com/JackalLabs/jackalgo/handlers/wallet_handler"

type StorageHandler struct {
	walletHandler *wallet_handler.WalletHandler
}

func NewStorageHandler(w *wallet_handler.WalletHandler) *StorageHandler {
	s := StorageHandler{
		walletHandler: w,
	}

	return &s
}
