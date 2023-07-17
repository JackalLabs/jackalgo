package compression

import (
	ecies "github.com/ecies/go/v2"
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
type BasePerms struct {
	TrackingNumber string
	Key            []byte
	Iv             []byte
}

type StandardPerms struct {
	BasePerms BasePerms
	PubKey    *ecies.PublicKey
	Usr       string
}

type EditorsViewers map[string]string
