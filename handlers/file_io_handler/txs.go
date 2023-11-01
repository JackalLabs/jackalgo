package file_io_handler

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/JackalLabs/jackalgo/handlers/storage_handler"
	"github.com/JackalLabs/jackalgo/utils"
	"strings"
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

func (f *FileIoHandler) SignAndBroadcast(msgs []sdk.Msg) error {
	res, err := f.walletHandler.SendTx(msgs...)
	if err != nil {
		return err
	}
	fmt.Println(res.Code)
	fmt.Println(res.RawLog)
	return nil
}

func (f *FileIoHandler) LoadNestedFolder(rawPath string) (folderHandlers *folder_handler.FolderHandler, msgs []sdk.Msg, err error) {
	folders := folder_handler.FolderGroup{}
	msgs = []sdk.Msg{}
	pathChunks := strings.Split(rawPath, "/")

	for i := 1; i < len(pathChunks); i++ {
		parentSubString := strings.Join(pathChunks[0:i], "/")
		subString := strings.Join(pathChunks[0:i+1], "/")

		fmt.Println("cycle start")
		fmt.Println(parentSubString)
		fmt.Println(subString)

		rawSubFolder, err := compression.ReadFileTreeEntry(f.walletHandler.GetAddress(), subString, f.walletHandler)
		fmt.Println(rawSubFolder)
		fmt.Println(err)
		fmt.Println(pathChunks[i])
		if err != nil {
			folders[subString] = folder_handler.TrackNewFolder(
				pathChunks[i],
				parentSubString,
				folders[parentSubString].GetWhoOwnsMe(),
				f.walletHandler,
			)

			fmt.Println(folders[subString].GetFolderDetails())

			msg, _, err := folders[parentSubString].AddChildDirs([]string{pathChunks[i]})
			fmt.Println(folders[parentSubString].GetChildDirs())
			fmt.Println(msg)
			fmt.Println(err)
			if err != nil {
				return nil, nil, err
			}

			msgs = append(msgs, msg...)
		} else {
			var subFrame folder_handler.FolderFileFrame
			err = json.Unmarshal(rawSubFolder, &subFrame)
			if err != nil {
				return nil, nil, err
			}

			folders[subString] = folder_handler.TrackFolder(subFrame, f.walletHandler)
		}
		fmt.Println("msgs")
		fmt.Println(msgs)
	}
	fmt.Println("LoadNestedFolder done")

	return folders[rawPath], msgs, err
}

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
		PubKey:    f.walletHandler.GetECIESPubKey(),
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

	//fmt.Println(dirMsgs)

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

func (f *FileIoHandler) StaggeredUploadFiles(sourceFiles []*file_upload_handler.FileUploadHandler, parent *folder_handler.FolderHandler, public bool) (failedCount int, fids []string, cids []string, err error) {
	failedFiles := make([]*types.File, 0)
	fids = make([]string, 0)
	cids = make([]string, 0)

	count := 0
	total := 0
	queue := NewQueue()
	counter := 0
	for _, file := range sourceFiles {
		outerFile := file
		counter++
		go func() {
			fmt.Printf("uploading %s...\n", outerFile.GetWhoAmI())
			defer func() {
				total++
				counter--
			}() // don't do this
			innerFile, err := outerFile.GetForUpload(public)
			if err != nil {
				failedFiles = append(failedFiles, innerFile)
				fmt.Printf("getting for upload failed for %s.\n", outerFile.GetWhoAmI())
				return
			}

			fid, cid, err := f.tumbleUpload(f.walletHandler.GetAddress(), innerFile)
			if err != nil {
				failedFiles = append(failedFiles, innerFile)
				fmt.Printf("failed to upload %s.\n", outerFile.GetWhoAmI())
				return
			}
			fids = append(fids, fid)
			cids = append(cids, cid)
			outerFile.SetIds(cid, []string{fid})
			fmt.Printf("done uploading %s with fid: %s.\n", outerFile.GetWhoAmI(), fid)
			queue.Push(outerFile)
			count++
		}()
	}

	for counter > 0 || count > 0 {
		fmt.Printf("waiting for files to finish...\n")
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
			fmt.Printf("Handling %s...\n", handler.GetWhoAmI())

			count--
			handlers = append(handlers, handler)

			metas[handler.GetMeta().Name] = handler.GetMeta()

		}
		msgs, err := f.signAndPostFiletree(handlers)
		if err != nil {
			continue
		}

		msg, err := parent.AddChildFileReferences(metas)
		if err != nil {
			return len(failedFiles) + failedCount, fids, cids, err
		}

		msgs = append(msgs, msg)

		res, err := f.walletHandler.SendTx(msgs...)
		if err != nil {
			failedCount += len(handlers)
		} else {
			fmt.Println(res.Code)
			fmt.Println(res.RawLog)
		}

	}

	return len(failedFiles) + failedCount, fids, cids, nil
}

func (f *FileIoHandler) signAndPostFiletree(handlers []*file_upload_handler.FileUploadHandler) ([]sdk.Msg, error) {
	toBroadcast := make([]sdk.Msg, 0)

	if len(handlers) == 0 {
		return nil, fmt.Errorf("no files to upload")
	}

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
			PubKey: f.walletHandler.GetECIESPubKey(),
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

func (f *FileIoHandler) DeleteTargets(targets []string, parent *folder_handler.FolderHandler) error {
	msgs, err := f.deleteTargets(targets, parent)
	if err != nil {
		return err
	}

	_, err = f.walletHandler.SendTx(msgs...)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileIoHandler) makeDelete(creator string, target string) ([]sdk.Msg, error) {
	s := storage_handler.NewStorageHandler(f.walletHandler)

	msgs := make([]sdk.Msg, 0)

	treeData, err := utils.GetFileTreeData(target, creator, f.walletHandler)
	if err != nil {
		return nil, err
	}

	fids := Fids{}

	err = json.Unmarshal([]byte(treeData.Files.Contents), &fids)
	if err != nil {
		return nil, err
	}

	cidRes, err := s.QueryFidCid(fids.Fids[0])
	if err != nil {
		return nil, err
	}

	var cids []string
	err = json.Unmarshal([]byte(cidRes.FidCid.Cids), &cids)
	if err != nil {
		return nil, err
	}

	for _, cid := range cids {
		msg := &storagetypes.MsgCancelContract{
			Creator: creator,
			Cid:     cid,
		}

		msgs = append(msgs, msg)
	}

	msg := &filetreetypes.MsgDeleteFile{
		Creator:  creator,
		HashPath: crypt.MerkleMeBro(target),
		Account:  crypt.HashAndHex(creator),
	}
	msgs = append(msgs, msg)

	return msgs, nil
}

func (f *FileIoHandler) deleteTargets(targets []string, parent *folder_handler.FolderHandler) ([]sdk.Msg, error) {
	childFiles := parent.GetChildFiles()
	childFolders := parent.GetChildDirs()

	msgs := make([]sdk.Msg, 0)
	location := fmt.Sprintf("%s/%s", parent.GetWhereAmI(), parent.GetWhoAmI())

	for _, childFile := range childFiles {
		for _, target := range targets {
			rawPath := fmt.Sprintf("%s/%s", location, target)
			if target == childFile.Name {
				deletionMessages, err := f.makeDelete(f.walletHandler.GetAddress(), rawPath)
				if err != nil {
					return nil, err
				}

				msg, err := parent.RemoveChildFileReferences([]string{target})
				if err != nil {
					return nil, err
				}
				msgs = append(msgs, msg)

				msgs = append(msgs, deletionMessages...)
			}
		}
	}
	for _, childFolder := range childFolders {
		for _, target := range targets {
			rawPath := fmt.Sprintf("%s/%s", location, target)
			if target == childFolder {
				innerFolder, err := f.DownloadFolder(rawPath)
				if err != nil {
					return nil, err
				}

				msg, err := parent.RemoveChildDirReferences([]string{target})
				if err != nil {
					return nil, err
				}
				msgs = append(msgs, msg)

				dirs := innerFolder.GetChildDirs()
				for _, file := range innerFolder.GetChildFiles() {
					dirs = append(dirs, file.Name)
				}

				deletionMessages, err := f.deleteTargets(dirs, innerFolder)
				if err != nil {
					return nil, err
				}

				msgs = append(msgs, deletionMessages...)
			}
		}
	}

	return msgs, nil
}
