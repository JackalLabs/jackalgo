package file_io_handler

import "github.com/JackalLabs/jackalgo/handlers/wallet_handler"

type FileIoHandler struct {
	walletHandler *wallet_handler.WalletHandler
}

func NewFileIoHandler(w *wallet_handler.WalletHandler) *FileIoHandler {
	f := FileIoHandler{
		walletHandler: w,
	}

	return &f
}
