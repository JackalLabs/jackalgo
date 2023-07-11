package file_io_handler

import (
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/JackalLabs/jackalgo/types"
	"strings"
)

type FileIoHandler struct {
	walletHandler *wallet_handler.WalletHandler
}

func NewFileIoHandler(w *wallet_handler.WalletHandler) *FileIoHandler {
	f := FileIoHandler{
		walletHandler: w,
	}

	return &f
}

func tumbleUpload(sender string, file *types.File) {
	// TODO - get provider
	prov := get()

	url := strings.Trim(prov, "/")
	fid, cid, err := doUpload(url, sender, file)
	if err != nil {
		// TODO - change provider and try again
	}
}

func doUpload(url string, sender string, file *types.File) (fid string, cid string, err error) {

	// TODO - doUpload
	// TODO - create FormData for receiving provider to read
	fileFormData.set('file', file)
	fileFormData.set('sender', sender)
	// TODO - create https upload process
	fid, cid, err := upload(url)
	if err != nil {
		return "", "", err
	}
	// TODO - end doUpload
	return fid, cid, nil
}
