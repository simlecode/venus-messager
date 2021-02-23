package impl

import (
	"context"

	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/db"
	"github.com/ipfs-force-community/venus-messager/lib/message"
	"github.com/ipfs-force-community/venus-messager/node/config"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/go-jsonrpc/auth"
	logging "github.com/ipfs/go-log/v2"

	"github.com/ipfs-force-community/venus-messager/api"
)

var log = logging.Logger("node")

type FullNodeAPI struct {
}

func (n *FullNodeAPI) MessagePush(ctx context.Context, msg *types.Message, spec *api.MessageSendSpec, meta *message.MsgMeta) error {
	panic("implement me")
}

func (n *FullNodeAPI) MessagesPush(ctx context.Context, msg []*types.Message, spec []*api.MessageSendSpec, meta []*message.MsgMeta) error {
	panic("implement me")
}

func (n *FullNodeAPI) MessagesPub(ctx context.Context, signedMsg []*types.SignedMessage) error {
	panic("implement me")
}

func (n *FullNodeAPI) NewDbProcess(config *config.SQLiteDBConfig) (db.DBProcessInterface, error) {
	panic("implement me")
}

func (n *FullNodeAPI) AuthVerify(ctx context.Context, token string) ([]auth.Permission, error) {
	panic("implement me")
}

func (n *FullNodeAPI) AuthNew(ctx context.Context, perms []auth.Permission) ([]byte, error) {
	panic("implement me")
}

func (n *FullNodeAPI) ID(context.Context) (peer.ID, error) {
	panic("implement me")
}

func (n *FullNodeAPI) Version(context.Context) (api.Version, error) {
	panic("implement me")
}

func (n *FullNodeAPI) LogList(context.Context) ([]string, error) {
	panic("implement me")
}

func (n *FullNodeAPI) LogSetLevel(context.Context, string, string) error {
	panic("implement me")
}

func (n *FullNodeAPI) Shutdown(context.Context) error {
	panic("implement me")
}

func (n *FullNodeAPI) Session(context.Context) (uuid.UUID, error) {
	panic("implement me")
}

func (n *FullNodeAPI) Closing(context.Context) (<-chan struct{}, error) {
	panic("implement me")
}

var _ api.FullNode = &FullNodeAPI{}
