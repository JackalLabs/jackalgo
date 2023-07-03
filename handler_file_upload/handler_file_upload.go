package handler_file_upload

import (
	"fmt"
)

func NewFileUploadHandler() *FileUploadHandler {

	f := FileUploadHandler{}

	return &f
}

func (f *FileUploadHandler) SayHello() {
	fmt.Println("Hello from FileUploadHandler")
}
