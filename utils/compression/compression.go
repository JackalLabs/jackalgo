package compression

import (
	"encoding/json"
	"fmt"
	"unicode/utf16"

	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils"
	lzstring "github.com/daku10/go-lz-string"
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

	msg := MsgPartialPostFileBundle{
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

	perms := basePerms{
		trackingNumber: msg.TrackingNumber,
		iv:             iv,
		key:            key,
	}

	selfPubKey := walletRef.GetPubKey()
	me := standardPerms{
		basePerms: perms,
		pubKey:    selfPubKey,
		usr:       creator,
	}

	ukey, uivkey, err := makePermsBlock("e", me, walletRef)
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
		ukey, uivkey, err := makePermsBlock("v", me, walletRef)
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
		them := standardPerms{
			basePerms: perms,
			pubKey:    destPubKey,
			usr:       toAddress,
		}

		ev := make(EditorsViewers, 0)
		r1key, r1ivkey, err := makePermsBlock("v", me, walletRef)
		if err != nil {
			return nil, err
		}
		r2key, r2ivkey, err := makePermsBlock("v", them, walletRef)
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

	return buildPostFile(msg), nil
}

func ReadFileTreeEntry(owner string, rawPath string, walletRef types.Wallet) (map[string]any, error) {
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

	var ff map[string]any
	err = json.Unmarshal([]byte(final), &ff)

	return ff, err
}

func buildPostFile(data MsgPartialPostFileBundle) sdk.Msg {
	return &filetreetypes.MsgPostFile{
		Creator:        data.Creator,
		Account:        data.Account,
		HashParent:     data.HashParent,
		HashChild:      data.HashChild,
		Contents:       data.Contents,
		Editors:        data.Editors,
		Viewers:        data.Viewers,
		TrackingNumber: data.TrackingNumber,
	}
}

func makePermsBlock(base string, standardPerms standardPerms, walletRef types.Wallet) (string, string, error) {
	user := crypt.HashAndHex(fmt.Sprintf("%s%s%s", base, standardPerms.basePerms.trackingNumber, standardPerms.usr))
	perms, err := crypt.AesToString(walletRef, standardPerms.pubKey, standardPerms.basePerms.key, standardPerms.basePerms.iv)
	if err != nil {
		return "", "", err
	}
	return user, perms, nil
}

func CompressData(input string) (string, error) {
	k, err := lzstring.Compress(input)
	if err != nil {
		return "", err
	}
	s := string(utf16.Decode(k))
	return fmt.Sprintf("jklpc1%s", s), nil
}

func DecompressData(input string) (string, error) {
	s := utf16.Encode([]rune(input))
	return lzstring.Decompress(s)
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
