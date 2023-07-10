package compression

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

type MsgPartialPostFileBundle struct {
	Creator        string // just the users address (might rework to be the same as account)
	Account        string // the hashed (uuid + user account) (becomes owner)
	HashParent     string // merkled parent
	HashChild      string // hashed child
	Contents       string // contents
	Viewers        string // stringify IEditorsViewers
	Editors        string // stringify IEditorsViewers
	TrackingNumber string // uuid
}
type basePerms struct {
	trackingNumber string
	key            []byte
	iv             []byte
}

type standardPerms struct {
	basePerms basePerms
	pubKey    cryptotypes.PubKey
	usr       string
}

type EditorsViewers map[string]string
