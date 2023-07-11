package file_io_handler

import (
	"github.com/JackalLabs/jackalgo/handlers/folder_handler"
	"github.com/JackalLabs/jackalgo/utils/compression"
)

func (f *FileIoHandler) DownloadFolder(rawPath string) (folderHandler *folder_handler.FolderHandler, err error) {
	rawFolder, err := compression.ReadFileTreeEntry(f.walletHandler.GetAddress(), rawPath, f.walletHandler)
	if err != nil {
		return nil, err
	}
	// TODO - convert generic ReadFileTreeEntry return to FolderFileFrame
	return folder_handler.TrackFolder(rawFolder, f.walletHandler), nil
}
