package full

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/google/uuid"
	"github.com/ipfs-force-community/venus-messager/api"
	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/message"
	"github.com/ipfs-force-community/venus-messager/node/config"
	"github.com/ipfs-force-community/venus-messager/node/impl/db"
	peer "github.com/libp2p/go-libp2p-peer"
	"go.uber.org/fx"
)

type MessageAPI struct {
	fx.In

	DB DBProcessAPI
}

func (m *MessageAPI) MessagePush(ctx context.Context, msg *types.Message, msgMeta *message.MsgMeta) error {
	return m.DB.AddMessage(msg, msgMeta)
}

func (m *MessageAPI) MessagesPush(ctx context.Context, msg []*types.Message, msgMeta []*message.MsgMeta) error {
	if len(msg) != len(msgMeta) {
		return xerrors.Errorf("message length(%d) not match message meta length(%d)", len(msg), len(msgMeta))
	}
	for i, msg := range msg {
		if err := m.DB.AddMessage(msg, msgMeta[i]); err != nil {
			return err
		}
	}

	return nil
}

func (m *MessageAPI) MessagesPub(ctx context.Context, signedMsg []*types.SignedMessage) error {
	panic("implement me")
}

func (m *MessageAPI) NewDbProcess(config *config.DBConfig) (db.DBProcessInterface, error) {
	panic("implement me")
}

func (m *MessageAPI) AuthVerify(ctx context.Context, token string) ([]auth.Permission, error) {
	panic("implement me")
}

func (m *MessageAPI) AuthNew(ctx context.Context, perms []auth.Permission) ([]byte, error) {
	panic("implement me")
}

func (m *MessageAPI) ID(context.Context) (peer.ID, error) {
	panic("implement me")
}

func (m *MessageAPI) Version(context.Context) (api.Version, error) {
	panic("implement me")
}

func (m *MessageAPI) LogList(context.Context) ([]string, error) {
	panic("implement me")
}

func (m *MessageAPI) LogSetLevel(context.Context, string, string) error {
	panic("implement me")
}

func (m *MessageAPI) Shutdown(context.Context) error {
	panic("implement me")
}

func (m *MessageAPI) Session(context.Context) (uuid.UUID, error) {
	panic("implement me")
}

func (m *MessageAPI) Closing(context.Context) (<-chan struct{}, error) {
	panic("implement me")
}
