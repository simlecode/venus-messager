package api

import (
	"context"

	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/message"
)

type FullNode interface {
	Common

	// message
	// MessagePush receiving messages from other nodes
	MessagePush(ctx context.Context, msg *types.Message, spec *MessageSendSpec, meta *message.MsgMeta) error
	MessagesPush(ctx context.Context, msg []*types.Message, spec []*MessageSendSpec, meta []*message.MsgMeta) error
	// MessagesPub Publish the signed messages
	MessagesPub(ctx context.Context, signedMsg []*types.SignedMessage) error

	// database process interface
}
