package file_download_handler

import (
	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils/crypt"
)

type FileDownloadHandler struct {
	File *types.File
}

func NewFileDownloadHandler(file *types.File) *FileDownloadHandler {
	f := FileDownloadHandler{
		File: file,
	}

	return &f
}

func TrackFile(file []byte, key []byte, iv []byte) (*FileDownloadHandler, error) {
	decryptedFile, err := crypt.ConvertFromEncryptedFile(file, key, iv)
	if err != nil {
		return nil, err
	}
	return NewFileDownloadHandler(decryptedFile), nil
}

func (f *FileDownloadHandler) GetFile() *types.File {
	return f.File
}
