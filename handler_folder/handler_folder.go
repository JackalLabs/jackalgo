package handler_folder

import (
	"fmt"
)

func NewFolderHandler() *FolderHandler {

	r := FolderHandler{}

	return &r
}

func (f *FolderHandler) SayHello() {
	fmt.Println("Hello from FolderHandler")
}
