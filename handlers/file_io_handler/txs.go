package file_io_handler

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JackalLabs/jackalgo/handlers/file_upload_handler"
	"github.com/JackalLabs/jackalgo/handlers/folder_handler"
	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils/compression"
	"github.com/JackalLabs/jackalgo/utils/crypt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	filetreetypes "github.com/jackalLabs/canine-chain/v3/x/filetree/types"
	storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"
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
		TrackingNumber: id,
		Iv:             crypt.GenIv(),
		Key:            crypt.GenKey(),
	}

	standard := compression.StandardPerms{
		BasePerms: base,
		PubKey:    f.walletHandler.GetPubKey(),
		Usr:       f.walletHandler.GetAddress(),
	}

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

	return &filetreetypes.MsgMakeRootV2{
		Creator:        f.walletHandler.GetAddress(),
		Editors:        string(editors),
		Viewers:        string(viewers),
		TrackingNumber: id,
	}, nil
}

func (f *FileIoHandler) GenerateInitialDirs(startingDirs []string) (*sdk.TxResponse, error) {
	toGenerate := startingDirs
	if len(toGenerate) == 0 {
		toGenerate = []string{"Config", "Home", "WWW"}
	}

	creator := f.walletHandler.GetAddress()

	dirMsgs := make([]sdk.Msg, len(toGenerate))
	for i, generation := range toGenerate {
		handler := folder_handler.TrackNewFolder(generation, "s", creator, f.walletHandler)
		msg, err := handler.GetForFiletree()
		if err != nil {
			return nil, err
		}
		dirMsgs[i] = msg
	}

	readyToBroadcast := make([]sdk.Msg, 0)

	pubKeyQuerier := filetreetypes.NewQueryClient(f.walletHandler.GetClientCtx())
	req := filetreetypes.QueryPubkeyRequest{
		Address: creator,
	}
	_, err := pubKeyQuerier.Pubkey(context.Background(), &req)
	if err != nil {
		hexKey := hex.EncodeToString(f.walletHandler.GetECIESPubKey().Bytes(false))

		postKey := filetreetypes.MsgPostkey{
			Creator: creator,
			Key:     hexKey,
		}
		readyToBroadcast = append(readyToBroadcast, &postKey)
	}

	root, err := f.CreateRoot()
	if err != nil {
		return nil, err
	}
	readyToBroadcast = append(readyToBroadcast, root)
	readyToBroadcast = append(readyToBroadcast, dirMsgs...)

	return f.walletHandler.SendTx(readyToBroadcast...)
}

func (f *FileIoHandler) StaggeredUploadFiles(sourceFiles []*file_upload_handler.FileUploadHandler, parent *folder_handler.FolderHandler, public bool) error {
	failedFiles := make([]*types.File, 0)

	count := 0
	total := 0
	queue := NewQueue()

	for _, file := range sourceFiles {
		outerFile := file
		go func() {
			defer func() { total++ }() // don't do this
			count++

			innerFile, err := outerFile.GetForUpload(public)
			if err != nil {
				failedFiles = append(failedFiles, innerFile)
				return
			}

			fid, cid, err := f.tumbleUpload(f.walletHandler.GetAddress(), innerFile)
			if err != nil {
				failedFiles = append(failedFiles, innerFile)
				return
			}
			outerFile.SetIds(cid, []string{fid})
			queue.Push(outerFile)
		}()
	}

	for count > 0 {
		for i := 0; i < 12; i++ {
			time.Sleep(5 * time.Second)
			if total == len(sourceFiles) {
				break
			}
		}

		metas := make(folder_handler.FolderChildFiles, 0)
		handlers := make([]*file_upload_handler.FileUploadHandler, 0)
		for !queue.Empty() {
			handler := queue.Pop()
			count--
			handlers = append(handlers, handler)

			metas[handler.GetMeta().Name] = handler.GetMeta()

		}
		msgs, err := f.signAndPostFiletree(handlers)
		if err != nil {
			return err
		}

		msg, err := parent.AddChildFileReferences(metas)
		if err != nil {
			return err
		}

		msgs = append(msgs, msg)

		_, err = f.walletHandler.SendTx(msgs...)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (f *FileIoHandler) signAndPostFiletree(handlers []*file_upload_handler.FileUploadHandler) ([]sdk.Msg, error) {
	toBroadcast := make([]sdk.Msg, 0)

	for _, handler := range handlers {
		cid, fids := handler.GetIds()

		fs := Fids{
			Fids: fids,
		}

		fidJson, err := json.Marshal(fs)
		if err != nil {
			return nil, err
		}

		key, iv := handler.GetEnc()
		perms := compression.StandardPerms{
			BasePerms: compression.BasePerms{
				TrackingNumber: handler.GetUUID(),
				Key:            key,
				Iv:             iv,
			},
			PubKey: f.walletHandler.GetPubKey(),
			Usr:    f.walletHandler.GetAddress(),
		}
		u, p, err := compression.MakePermsBlock("v", perms, f.walletHandler)
		if err != nil {
			return nil, err
		}
		vev := make(compression.EditorsViewers, 0)
		vev[u] = p
		jsonViewers, err := json.Marshal(vev)
		if err != nil {
			return nil, err
		}

		u, p, err = compression.MakePermsBlock("e", perms, f.walletHandler)
		if err != nil {
			return nil, err
		}
		eev := make(compression.EditorsViewers, 0)
		eev[u] = p
		jsonEditors, err := json.Marshal(vev)
		if err != nil {
			return nil, err
		}

		msgPostFileBundle := filetreetypes.MsgPostFile{
			Creator:        f.walletHandler.GetAddress(),
			Account:        crypt.HashAndHex(f.walletHandler.GetAddress()),
			HashParent:     handler.GetMerklePath(),
			HashChild:      crypt.HashAndHex(handler.GetWhoAmI()),
			Contents:       string(fidJson),
			Viewers:        string(jsonViewers),
			Editors:        string(jsonEditors),
			TrackingNumber: handler.GetUUID(),
		}

		msgSign := storagetypes.MsgSignContract{
			Creator: f.walletHandler.GetAddress(),
			Cid:     cid,
			PayOnce: false,
		}

		toBroadcast = append(toBroadcast, &msgSign)
		toBroadcast = append(toBroadcast, &msgPostFileBundle)

	}

	return toBroadcast, nil
}
