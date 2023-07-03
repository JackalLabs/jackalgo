package handler_file_io

import (
	"fmt"
	"github.com/JackalLabs/jackalgo"
)

func NewFileIoHandler(w *jackalgo.WalletHandler) *FileIoHandler {

	f := FileIoHandler{
		walletHandler: w,
	}

	return &f
}

func (f *FileIoHandler) SayHello() {
	fmt.Println("Hello from FileIoHandler")
}
