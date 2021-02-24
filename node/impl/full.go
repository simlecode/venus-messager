package impl

import (
	logging "github.com/ipfs/go-log/v2"

	"github.com/ipfs-force-community/venus-messager/api"
	"github.com/ipfs-force-community/venus-messager/node/impl/full"
)

var log = logging.Logger("node")

type FullNodeAPI struct {
	full.DBProcessAPI
	full.MessageAPI
}

var _ api.FullNode = &FullNodeAPI{}
