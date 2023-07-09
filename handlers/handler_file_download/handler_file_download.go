package handler_file_download

import (
	"fmt"
)

func NewFileDownloadHandler() *FileDownloadHandler {

	f := FileDownloadHandler{}

	return &f
}

func (f *FileDownloadHandler) SayHello() {
	fmt.Println("Hello from FileDownloadHandler")
}
