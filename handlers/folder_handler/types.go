package folder_handler

import (
	"fmt"

	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils/compression"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FolderHandler struct {
	folderDetails FolderFileFrame
	isFolder      bool
	walletHandler *wallet_handler.WalletHandler
}

func NewFolderHandler(frame FolderFileFrame, wallet *wallet_handler.WalletHandler) *FolderHandler {
	r := FolderHandler{
		folderDetails: frame,
		isFolder:      true,
		walletHandler: wallet,
	}

	return &r
}

func TrackFolder(dirInfo FolderFileFrame, wallet *wallet_handler.WalletHandler) *FolderHandler {
	return NewFolderHandler(dirInfo, wallet)
}

func (f *FolderHandler) GetWhoAmI() string {
	return f.folderDetails.WhoAmI
}

func (f *FolderHandler) GetWhereAmI() string {
	return f.folderDetails.WhereAmI
}

func (f *FolderHandler) GetWhoOwnsMe() string {
	return f.folderDetails.WhoOwnsMe
}

func (f *FolderHandler) GetMyPath() string {
	return fmt.Sprintf("%s/%s", f.folderDetails.WhereAmI, f.folderDetails.WhoAmI)
}

func (f *FolderHandler) GetMyChildPath(child string) string {
	return fmt.Sprintf("%s/%s", f.GetMyPath(), child)
}

func (f *FolderHandler) GetFolderDetails() FolderFileFrame {
	return f.folderDetails
}

func (f *FolderHandler) GetChildDirs() []string {
	return f.folderDetails.DirChildren
}

func (f *FolderHandler) GetChildFiles() FolderChildFiles {
	return f.folderDetails.FileChildren
}

func (f *FolderHandler) GetForFiletree() (sdk.Msg, error) {
	return compression.SaveFiletreeEntry(f.walletHandler.GetAddress(), f.GetWhereAmI(), f.GetWhoAmI(), f.folderDetails, f.walletHandler)
}

type FolderFileFrame struct {
	WhoAmI       string
	WhereAmI     string
	WhoOwnsMe    string
	DirChildren  []string
	FileChildren FolderChildFiles
}

type FolderChildFiles map[string]types.Details
