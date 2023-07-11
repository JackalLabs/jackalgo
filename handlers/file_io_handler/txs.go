package file_io_handler

import (
	"encoding/json"
	"fmt"
	"github.com/JackalLabs/jackalgo/handlers/folder_handler"
	"github.com/JackalLabs/jackalgo/utils/compression"
	"github.com/JackalLabs/jackalgo/utils/crypt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	filetreetypes "github.com/jackalLabs/canine-chain/v3/x/filetree/types"
)

func (f *FileIoHandler) CreateFolders(parentDir *folder_handler.FolderHandler, newDirs []string) (msgs []sdk.Msg, err error) {
	msgs, existing, err := parentDir.AddChildDirs(newDirs)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		fmt.Printf("The following folders already exist: %s", existing)
	}
	return msgs, nil
}

func (f *FileIoHandler) CreateRoot() (msgs sdk.Msg, err error) {
	id := uuid.New().String()

	base := compression.BasePerms{
		trackingNumber: id,
		iv:             crypt.GenIv(),
		key:            crypt.GenKey(),
	}

	standard := compression.StandardPerms{
		basePerms: base,
		pubKey:    f.walletHandler.GetPubKey(),
		usr:       f.walletHandler.GetAddress(),
	}

	compression.MakePermsBlock("e", standard, f.walletHandler)

	eukey, euivkey, err := compression.MakePermsBlock("e", standard, f.walletHandler)
	if err != nil {
		return nil, err
	}
	eev := make(compression.EditorsViewers, 0)
	eev[eukey] = euivkey

	editors, err := json.Marshal(eev)
	if err != nil {
		return nil, err
	}

	vukey, vuivkey, err := compression.MakePermsBlock("v", standard, f.walletHandler)
	if err != nil {
		return nil, err
	}
	vev := make(compression.EditorsViewers, 0)
	vev[vukey] = vuivkey

	viewers, err := json.Marshal(vev)
	if err != nil {
		return nil, err
	}

	// TODO - why is this throwing an error
	return &filetreetypes.MsgMakeRootV2{
		Creator:        f.walletHandler.GetAddress(),
		Editors:        string(editors),
		Viewers:        string(viewers),
		TrackingNumber: id,
	}

}
