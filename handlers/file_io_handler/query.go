package file_io_handler

import (
	"encoding/json"

	"github.com/JackalLabs/jackalgo/handlers/folder_handler"
	"github.com/JackalLabs/jackalgo/utils/compression"
)

func (f *FileIoHandler) DownloadFolder(rawPath string) (folderHandler *folder_handler.FolderHandler, err error) {
	rawFolder, err := compression.ReadFileTreeEntry(f.walletHandler.GetAddress(), rawPath, f.walletHandler)
	if err != nil {
		return nil, err
	}

	var frame folder_handler.FolderFileFrame
	err = json.Unmarshal(rawFolder, &frame)
	if err != nil {
		return nil, err
	}

	return folder_handler.TrackFolder(frame, f.walletHandler), nil
}
