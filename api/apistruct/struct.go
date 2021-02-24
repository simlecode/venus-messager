package apistruct

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/ipfs-force-community/venus-messager/api"
	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/message"
	"github.com/ipfs-force-community/venus-messager/node/impl/db"
)

// All permissions are listed in permissioned.go
var _ = AllPermissions

type CommonStruct struct {
	Internal struct {
		AuthVerify func(ctx context.Context, token string) ([]auth.Permission, error) `perm:"read"`
		AuthNew    func(ctx context.Context, perms []auth.Permission) ([]byte, error) `perm:"admin"`

		ID      func(context.Context) (peer.ID, error)     `perm:"read"`
		Version func(context.Context) (api.Version, error) `perm:"read"`

		LogList     func(context.Context) ([]string, error)     `perm:"write"`
		LogSetLevel func(context.Context, string, string) error `perm:"write"`
	}
}

// FullNodeStruct implements API passing calls to user-provided function values.
type FullNodeStruct struct {
	CommonStruct

	Internal struct {
		MessagePush func(ctx context.Context, msg *types.Message, meta *message.MsgMeta) error `perm:"write"`

		QueryMessage        func(addr address.Address, from, to uint64) ([]db.Msg, error) `perm:"read"`
		DelMessage          func(addr address.Address, from, to uint64) error             `perm:"write"`
		AddMessage          func(msg *types.Message, msgMeta *message.MsgMeta) error      `perm:"write"`
		UpdateSignedMessage func(id uint64, signedMsg *types.SignedMessage) error         `perm:"write"`
		SetNonce            func(addr address.Address, nonce uint64) error                `perm:"write"`
		QueryNonce          func(addr address.Address) (uint64, error)                    `perm:"write"`
	}
}

// CommonStruct

func (c *CommonStruct) AuthVerify(ctx context.Context, token string) ([]auth.Permission, error) {
	return c.Internal.AuthVerify(ctx, token)
}

func (c *CommonStruct) AuthNew(ctx context.Context, perms []auth.Permission) ([]byte, error) {
	return c.Internal.AuthNew(ctx, perms)
}

// ID implements API.ID
func (c *CommonStruct) ID(ctx context.Context) (peer.ID, error) {
	return c.Internal.ID(ctx)
}

// Version implements API.Version
func (c *CommonStruct) Version(ctx context.Context) (api.Version, error) {
	return c.Internal.Version(ctx)
}

func (c *CommonStruct) LogList(ctx context.Context) ([]string, error) {
	return c.Internal.LogList(ctx)
}

func (c *CommonStruct) LogSetLevel(ctx context.Context, group, level string) error {
	return c.Internal.LogSetLevel(ctx, group, level)
}

// FullNodeStruct

func (f FullNodeStruct) MessagePush(ctx context.Context, msg *types.Message, meta *message.MsgMeta) error {
	return f.Internal.MessagePush(ctx, msg, meta)
}

func (f FullNodeStruct) MessagesPush(ctx context.Context, msg []*types.Message, meta []*message.MsgMeta) error {
	panic("implement me")
}

func (f FullNodeStruct) MessagesPub(ctx context.Context, signedMsg []*types.SignedMessage) error {
	panic("implement me")
}

func (f FullNodeStruct) QueryMessage(addr address.Address, from, to uint64) ([]db.Msg, error) {
	return f.Internal.QueryMessage(addr, from, to)
}

func (f FullNodeStruct) DelMessage(addr address.Address, from, to uint64) error {
	return f.Internal.DelMessage(addr, from, to)
}

func (f FullNodeStruct) AddMessage(msg *types.Message, msgMeta *message.MsgMeta) error {
	return f.Internal.AddMessage(msg, msgMeta)
}

func (f FullNodeStruct) UpdateSignedMessage(id uint64, signedMsg *types.SignedMessage) error {
	return f.Internal.UpdateSignedMessage(id, signedMsg)
}

func (f FullNodeStruct) SetNonce(addr address.Address, nonce uint64) error {
	return f.Internal.SetNonce(addr, nonce)
}

func (f FullNodeStruct) QueryNonce(addr address.Address) (uint64, error) {
	return f.Internal.QueryNonce(addr)
}

var _ api.Common = &CommonStruct{}
var _ api.FullNode = &FullNodeStruct{}
