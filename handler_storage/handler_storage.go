package handler_storage

import (
	"fmt"
	"github.com/JackalLabs/jackalgo"
)

func NewStorageHandler(w *jackalgo.WalletHandler) *StorageHandler {

	s := StorageHandler{
		walletHandler: w,
	}

	return &s
}

func (s *StorageHandler) SayHello() {
	fmt.Println("Hello from RnsHandler")
}
