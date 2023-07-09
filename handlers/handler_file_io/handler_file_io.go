package handler_file_io

import (
	"fmt"

	"github.com/JackalLabs/jackalgo/handlers/handler_wallet"
)

func NewFileIoHandler(w *handler_wallet.WalletHandler) *FileIoHandler {

	f := FileIoHandler{
		walletHandler: w,
	}

	return &f
}

func (f *FileIoHandler) SayHello() {
	fmt.Println("Hello from FileIoHandler")
}
