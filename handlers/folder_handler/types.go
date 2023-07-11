package folder_handler

import (
	"fmt"

	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils/compression"
	"github.com/JackalLabs/jackalgo/utils/crypt"
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

func TrackNewFolder(myName string, myParent string, myOwner string, wallet *wallet_handler.WalletHandler) *FolderHandler {
	folderDetails := FolderFileFrame{
		WhoAmI:       myName, // TODO : sanitize input
		WhereAmI:     myParent,
		WhoOwnsMe:    myOwner,
		DirChildren:  make([]string, 0),
		FileChildren: make(FolderChildFiles, 0),
	}

	return NewFolderHandler(folderDetails, wallet)
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

func (f *FolderHandler) GetChildMerkle(child string) string {
	return crypt.MerkleMeBro(fmt.Sprintf("%s/%s/%s", f.GetWhereAmI(), f.GetWhoAmI(), child))
}

func (f *FolderHandler) AddChildDirs(childNames []string) ([]sdk.Msg, []string, error) {
	existing := make([]string, 0)
	more := make([]string, 0)

	for _, name := range childNames {
		for _, childDir := range f.folderDetails.DirChildren {
			if name == childDir {
				existing = append(existing, name)
				continue
			}
			more = append(more, name)
		}
	}

	msgs := make([]sdk.Msg, len(more))
	for i, moreName := range more {
		myName, myPath, myOwner := f.MakeChildDirInfo(moreName)
		tracker := TrackNewFolder(myName, myPath, myOwner, f.walletHandler)
		ftree, err := tracker.GetForFiletree()
		if err != nil {
			return nil, nil, err
		}
		msgs[i] = ftree
	}

	if len(more) > 0 {
		set := make(map[string]interface{})
		for _, moreName := range more {
			set[moreName] = nil
		}
		for _, dirChild := range f.folderDetails.DirChildren {
			set[dirChild] = nil
		}

		children := make([]string, 0)
		for key := range set {
			children = append(children, key)
		}

		f.folderDetails.DirChildren = children
		filetreeMsg, err := f.GetForFiletree()
		if err != nil {
			return nil, nil, err
		}
		msgs = append(msgs, filetreeMsg)
	}

	return msgs, existing, nil
}

func (f *FolderHandler) AddChildFileReferences(newFiles FolderChildFiles) (sdk.Msg, error) {
	for key, value := range newFiles {
		f.folderDetails.FileChildren[key] = value
	}

	return f.GetForFiletree()
}

func (f *FolderHandler) RemoveChildDirReferences(toRemove []string) (sdk.Msg, error) {
	filtered := make([]string, 0)
	for _, child := range f.folderDetails.DirChildren {
		needsRemoved := false
		for _, removing := range toRemove {
			if child == removing {
				needsRemoved = true
				break
			}
		}
		if !needsRemoved {
			filtered = append(filtered, child)
		}
	}

	f.folderDetails.DirChildren = filtered

	return f.GetForFiletree()
}

func (f *FolderHandler) RemoveChildFileReferences(toRemove []string) (sdk.Msg, error) {
	for _, removing := range toRemove {
		delete(f.folderDetails.FileChildren, removing)
	}

	return f.GetForFiletree()
}

func (f *FolderHandler) MakeChildDirInfo(childName string) (myName string, myParent string, myOwner string) {
	myName = childName
	myParent = f.GetMyPath()
	myOwner = f.folderDetails.WhoOwnsMe
	return
}

type FolderFileFrame struct {
	WhoAmI       string           `json:"whoAmI"`
	WhereAmI     string           `json:"whereAmI"`
	WhoOwnsMe    string           `json:"whoOwnsMe"`
	DirChildren  []string         `json:"dirChildren"`
	FileChildren FolderChildFiles `json:"fileChildren"`
}

type FolderChildFiles map[string]types.Details
