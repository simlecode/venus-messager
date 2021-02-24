package db

import (
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/venus-messager/lib/message"

	"github.com/ipfs-force-community/venus-messager/chain/types"
)

type DBProcessInterface interface {
	QueryMessage(addr address.Address, from, to uint64) ([]Msg, error)
	DelMessage(addr address.Address, from, to uint64) error
	AddMessage(msg *types.Message, msgMeta *message.MsgMeta) error
	UpdateSignedMessage(id uint64, signedMsg *types.SignedMessage) error

	SetNonce(addr address.Address, nonce uint64) error
	QueryNonce(addr address.Address) (uint64, error)
}
