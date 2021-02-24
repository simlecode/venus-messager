package full

import (
	"go.uber.org/fx"

	"github.com/ipfs-force-community/venus-messager/node/impl/db"
)

type DBProcessAPI struct {
	fx.In

	db.DBProcessInterface
}
