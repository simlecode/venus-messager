package modules

import (
	"go.uber.org/fx"

	"github.com/ipfs-force-community/venus-messager/build"
	"github.com/ipfs-force-community/venus-messager/lib/peermgr"
	"github.com/ipfs-force-community/venus-messager/node/modules/dtypes"
	"github.com/ipfs-force-community/venus-messager/node/modules/helpers"
)

func RunPeerMgr(mctx helpers.MetricsCtx, lc fx.Lifecycle, pmgr *peermgr.PeerMgr) {
	go pmgr.Run(helpers.LifecycleCtx(mctx, lc))
}

func BuiltinDrandConfig() dtypes.DrandConfig {
	return build.DrandConfig
}
