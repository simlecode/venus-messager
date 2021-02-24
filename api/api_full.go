package api

import (
	"context"

	"github.com/filecoin-project/go-address"

	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/message"
	"github.com/ipfs-force-community/venus-messager/node/impl/db"
)

type FullNode interface {
	Common

	// message
	// MessagePush receiving messages from other nodes
	MessagePush(ctx context.Context, msg *types.Message, meta *message.MsgMeta) error
	MessagesPush(ctx context.Context, msg []*types.Message, meta []*message.MsgMeta) error
	// MessagesPub Publish the signed messages
	MessagesPub(ctx context.Context, signedMsg []*types.SignedMessage) error

	// database process interface
	QueryMessage(addr address.Address, from, to uint64) ([]db.Msg, error)
	DelMessage(addr address.Address, from, to uint64) error
	AddMessage(msg *types.Message, msgMeta *message.MsgMeta) error
	UpdateSignedMessage(id uint64, signedMsg *types.SignedMessage) error
	SetNonce(addr address.Address, nonce uint64) error
	QueryNonce(addr address.Address) (uint64, error)
}
