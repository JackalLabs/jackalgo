package compression

import (
	"encoding/json"
	"fmt"

	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils"
	lzstring "github.com/Lazarus/lz-string-go"
	filetreetypes "github.com/jackalLabs/canine-chain/v3/x/filetree/types"

	"github.com/JackalLabs/jackalgo/utils/crypt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

func SaveFiletreeEntry(toAddress string, rawPath string, rawTarget string, rawContents any, walletRef types.Wallet) (sdk.Msg, error) {
	creator := walletRef.GetAddress()
	account := crypt.HashAndHex(creator)

	iv := crypt.GenIv()
	key := crypt.GenKey()

	msg := filetreetypes.MsgPostFile{
		Account:        account,
		Creator:        creator,
		Contents:       "",
		HashParent:     crypt.MerkleMeBro(rawPath),
		HashChild:      crypt.HashAndHex(rawTarget),
		TrackingNumber: uuid.New().String(),
		Editors:        "",
		Viewers:        "",
	}
	jsonContents, err := json.Marshal(rawContents)
	if err != nil {
		return nil, err
	}

	msg.Contents, err = CompressEncryptString(
		string(jsonContents),
		key,
		iv,
	)
	if err != nil {
		return nil, err
	}

	perms := BasePerms{
		TrackingNumber: msg.TrackingNumber,
		Iv:             iv,
		Key:            key,
	}

	selfPubKey := walletRef.GetECIESPubKey()
	me := StandardPerms{
		BasePerms: perms,
		PubKey:    selfPubKey,
		Usr:       creator,
	}

	ukey, uivkey, err := MakePermsBlock("e", me, walletRef)
	if err != nil {
		return nil, err
	}
	ev := make(EditorsViewers, 0)
	ev[ukey] = uivkey

	editors, err := json.Marshal(ev)
	if err != nil {
		return nil, err
	}
	msg.Editors = string(editors)

	if toAddress == creator {
		ukey, uivkey, err := MakePermsBlock("v", me, walletRef)
		if err != nil {
			return nil, err
		}
		ev := make(EditorsViewers, 0)
		ev[ukey] = uivkey
		viewers, err := json.Marshal(ev)
		if err != nil {
			return nil, err
		}
		msg.Viewers = string(viewers)

	} else {
		destPubKey, err := walletRef.FindPubKey(toAddress)
		if err != nil {
			return nil, err
		}
		them := StandardPerms{
			BasePerms: perms,
			PubKey:    destPubKey,
			Usr:       toAddress,
		}

		ev := make(EditorsViewers, 0)
		r1key, r1ivkey, err := MakePermsBlock("v", me, walletRef)
		if err != nil {
			return nil, err
		}
		r2key, r2ivkey, err := MakePermsBlock("v", them, walletRef)
		if err != nil {
			return nil, err
		}
		ev[r1key] = r1ivkey
		ev[r2key] = r2ivkey

		viewers, err := json.Marshal(ev)
		if err != nil {
			return nil, err
		}
		msg.Viewers = string(viewers)
	}

	return &msg, nil
}

func ReadFileTreeEntry(owner string, rawPath string, walletRef types.Wallet) ([]byte, error) {
	result, err := utils.GetFileTreeData(rawPath, owner, walletRef)
	if err != nil {
		return nil, err
	}

	contents := result.Files.Contents
	viewingAccess := result.Files.ViewingAccess
	trackingNumber := result.Files.TrackingNumber
	var parsedViewingAccess EditorsViewers
	err = json.Unmarshal([]byte(viewingAccess), &parsedViewingAccess)
	if err != nil {
		return nil, err
	}

	viewName := crypt.HashAndHex(fmt.Sprintf("v%s%s", trackingNumber, walletRef.GetAddress()))

	iv, key, err := crypt.StringToAes(walletRef, parsedViewingAccess[viewName])
	if err != nil {
		return nil, err
	}

	final, err := DecryptDecompressString(contents, key, iv)
	if err != nil {
		return nil, err
	}

	return []byte(final), err
}

func MakePermsBlock(base string, standardPerms StandardPerms, walletRef types.Wallet) (user string, perms string, err error) {
	user = crypt.HashAndHex(fmt.Sprintf("%s%s%s", base, standardPerms.BasePerms.TrackingNumber, standardPerms.Usr))
	perms, err = crypt.AesToString(walletRef, standardPerms.PubKey, standardPerms.BasePerms.Key, standardPerms.BasePerms.Iv)
	if err != nil {
		return "", "", err
	}
	return
}

func CompressData(input string) (string, error) {
	k := lzstring.Compress(input, "")

	return fmt.Sprintf("jklpc1%s", k), nil
}

func DecompressData(input string) (string, error) {
	input = input[6:]
	return lzstring.Decompress(input, "")
}

func CompressEncryptString(input string, key []byte, iv []byte) (string, error) {
	compString, err := CompressData(input)
	if err != nil {
		return "", err
	}
	data, err := crypt.Encrypt([]byte(compString), key, iv)
	return string(data), err
}

func DecryptDecompressString(input string, key []byte, iv []byte) (string, error) {
	data, err := crypt.Decrypt([]byte(input), key, iv)
	if err != nil {
		return "", err
	}

	decompString, err := DecompressData(string(data))

	return decompString, err
}
