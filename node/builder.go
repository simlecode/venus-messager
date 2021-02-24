package node

import (
	"context"
	"errors"
	"time"

	logging "github.com/ipfs/go-log"
	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/fx"
	"golang.org/x/xerrors"

	"github.com/ipfs-force-community/venus-messager/api"
	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/peermgr"
	"github.com/ipfs-force-community/venus-messager/node/config"
	"github.com/ipfs-force-community/venus-messager/node/impl"
	"github.com/ipfs-force-community/venus-messager/node/impl/db"
	"github.com/ipfs-force-community/venus-messager/node/modules"
	"github.com/ipfs-force-community/venus-messager/node/modules/dtypes"
	"github.com/ipfs-force-community/venus-messager/node/modules/helpers"
	"github.com/ipfs-force-community/venus-messager/node/modules/lp2p"
	"github.com/ipfs-force-community/venus-messager/node/repo"
)

var log = logging.Logger("builder")

// special is a type used to give keys to modules which
//  can't really be identified by the returned type
type special struct{ id int }

// nolint:golint
var (
	DefaultTransportsKey = special{0}  // Libp2p option
	DiscoveryHandlerKey  = special{2}  // Private type
	AddrsFactoryKey      = special{3}  // Libp2p option
	SmuxTransportKey     = special{4}  // Libp2p option
	RelayKey             = special{5}  // Libp2p option
	SecurityKey          = special{6}  // Libp2p option
	BaseRoutingKey       = special{7}  // fx groups + multiret
	NatPortMapKey        = special{8}  // Libp2p option
	ConnectionManagerKey = special{9}  // Libp2p option
	AutoNATSvcKey        = special{10} // Libp2p option
)

type invoke int

// nolint:golint
const (
	// libp2p

	PstoreAddSelfKeysKey = invoke(iota)
	StartListeningKey
	BootstrapKey

	// filecoin
	SetGenesisKey

	RunHelloKey
	RunBlockSyncKey
	RunChainGraphsync
	RunPeerMgrKey

	HandleIncomingBlocksKey
	HandleIncomingMessagesKey

	RunDealClientKey
	RegisterClientValidatorKey

	// storage miner
	GetParamsKey
	HandleDealsKey
	HandleRetrievalKey
	RunSectorServiceKey
	RegisterProviderValidatorKey

	// daemon
	ExtractApiKey
	HeadMetricsKey
	RunPeerTaggerKey

	SetApiEndpointKey

	HandleMsgKey // 定时pub未上链消息

	_nInvokes // keep this last
)

type Settings struct {
	// modules is a map of constructors for DI
	//
	// In most cases the index will be a reflect. Type of element returned by
	// the constructor, but for some 'constructors' it's hard to specify what's
	// the return type should be (or the constructor returns fx group)
	modules map[interface{}]fx.Option

	// invokes are separate from modules as they can't be referenced by return
	// type, and must be applied in correct order
	invokes []fx.Option

	nodeType repo.RepoType

	Online bool // Online option applied
	Config bool // Config option applied

	Local bool
}

func defaults() []Option {
	return []Option{
		Override(new(helpers.MetricsCtx), context.Background),
		Override(new(record.Validator), modules.RecordValidator),
		Override(new(dtypes.Bootstrapper), dtypes.Bootstrapper(false)),
	}
}

func libp2p() Option {
	return Options(
		Override(new(peerstore.Peerstore), pstoremem.NewPeerstore),

		Override(DefaultTransportsKey, lp2p.DefaultTransports),

		Override(new(lp2p.RawHost), lp2p.Host),
		Override(new(host.Host), lp2p.RoutedHost),
		Override(new(lp2p.BaseIpfsRouting), lp2p.DHTRouting(dht.ModeAuto)),

		Override(DiscoveryHandlerKey, lp2p.DiscoveryHandler),
		Override(AddrsFactoryKey, lp2p.AddrsFactory(nil, nil)),
		Override(SmuxTransportKey, lp2p.SmuxTransport(true)),
		Override(RelayKey, lp2p.NoRelay()),
		Override(SecurityKey, lp2p.Security(true, true)),

		Override(BaseRoutingKey, lp2p.BaseRouting),
		Override(new(routing.Routing), lp2p.Routing),

		Override(NatPortMapKey, lp2p.NatPortMap),

		Override(ConnectionManagerKey, lp2p.ConnectionManager(50, 200, 20*time.Second, nil)),
		Override(AutoNATSvcKey, lp2p.AutoNATService),

		Override(new(*dtypes.ScoreKeeper), lp2p.ScoreKeeper),
		Override(new(*pubsub.PubSub), lp2p.GossipSub),
		Override(new(*config.Pubsub), func(bs dtypes.Bootstrapper) *config.Pubsub {
			return &config.Pubsub{
				Bootstrapper: bool(bs),
			}
		}),
		Override(PstoreAddSelfKeysKey, lp2p.PstoreAddSelfKeys),
		Override(StartListeningKey, lp2p.StartListening(config.DefaultFullNode().Libp2p.ListenAddresses)),
	)
}

func isType(t repo.RepoType) func(s *Settings) bool {
	return func(s *Settings) bool { return s.nodeType == t }
}

func Local() Option {
	return Options(
		// make sure that local is applied before Config.
		// This is important ***
		func(s *Settings) error { s.Local = true; return nil },
	)
}

// Online sets up basic libp2p node
func Online() Option {
	return Options(
		// make sure that online is applied before Config.
		// This is important because Config overrides some of Online units
		func(s *Settings) error { s.Online = true; return nil },
		ApplyIf(func(s *Settings) bool { return s.Config },
			Error(errors.New("the Online option must be set before Config option")),
		),

		libp2p(),

		// common

		// Full node
		ApplyIf(isType(repo.FullNode),
			Override(new(*peermgr.PeerMgr), peermgr.NewPeerMgr),
			Override(new(dtypes.DrandBootstrap), modules.DrandBootstrap),
			Override(new(dtypes.DrandConfig), modules.BuiltinDrandConfig),

			Override(new(dtypes.NetworkName), NetworkName),
			Override(new(db.DBProcessInterface), newDbProc),

			Override(RunPeerMgrKey, modules.RunPeerMgr),
		),
	)
}

// Config sets up constructors based on the provided Config
func ConfigCommon(cfg *config.Common) Option {
	return Options(
		func(s *Settings) error { s.Config = true; return nil },
		Override(new(dtypes.APIEndpoint), func() (dtypes.APIEndpoint, error) {
			return multiaddr.NewMultiaddr(cfg.API.ListenAddress)
		}),
		Override(SetApiEndpointKey, func(lr repo.LockedRepo, e dtypes.APIEndpoint) error {
			return lr.SetAPIEndpoint(e)
		}),
		ApplyIf(func(s *Settings) bool { return s.Online },
			Override(new(*pubsub.PubSub), lp2p.GossipSub),
			Override(new(*config.Pubsub), &cfg.Pubsub),
			Override(new(*config.DBConfig), &cfg.DBCfg),
			Override(new(*config.RPCServer), &cfg.RPCServer),
			ApplyIf(func(s *Settings) bool { return s.Local },
				Override(new(dtypes.BootstrapPeers), modules.ConfigBootstrap(cfg.Libp2p.BootstrapPeers, true)),
			),
			ApplyIf(func(s *Settings) bool { return !s.Local },
				Override(new(dtypes.BootstrapPeers), modules.ConfigBootstrap(cfg.Libp2p.BootstrapPeers, false)),
			),
		),
	)
}

func ConfigFullNode(c interface{}) Option {
	cfg, ok := c.(*config.FullNode)
	if !ok {
		return Error(xerrors.Errorf("invalid config from repo, got: %T", c))
	}

	return Options(
		ConfigCommon(&cfg.Common),
	)
}

func Repo(r repo.Repo) Option {
	return func(settings *Settings) error {
		lr, err := r.Lock(settings.nodeType)
		if err != nil {
			return err
		}
		c, err := lr.Config()
		if err != nil {
			return err
		}

		return Options(
			Override(new(repo.LockedRepo), modules.LockedRepo(lr)), // module handles closing

			Override(new(dtypes.MetadataDS), modules.Datastore),

			Override(new(ci.PrivKey), lp2p.PrivKey),
			Override(new(ci.PubKey), ci.PrivKey.GetPublic),
			Override(new(peer.ID), peer.IDFromPublicKey),

			Override(new(types.KeyStore), modules.KeyStore),

			Override(new(*dtypes.APIAlg), modules.APISecret),

			ApplyIf(isType(repo.FullNode), ConfigFullNode(c)),
		)(settings)
	}
}

func FullAPI(out *api.FullNode) Option {
	return func(s *Settings) error {
		resAPI := &impl.FullNodeAPI{}
		s.invokes[ExtractApiKey] = fx.Extract(resAPI)
		*out = resAPI
		return nil
	}
}

type StopFunc func(context.Context) error

// New builds and starts new Filecoin node
func New(ctx context.Context, opts ...Option) (StopFunc, error) {
	settings := Settings{
		modules:  map[interface{}]fx.Option{},
		invokes:  make([]fx.Option, _nInvokes),
		nodeType: repo.FullNode,
	}

	// apply module options in the right order
	if err := Options(Options(defaults()...), Options(opts...))(&settings); err != nil {
		return nil, xerrors.Errorf("applying node options failed: %w", err)
	}

	// gather constructors for fx.Options
	ctors := make([]fx.Option, 0, len(settings.modules))
	for _, opt := range settings.modules {
		ctors = append(ctors, opt)
	}

	// fill holes in invokes for use in fx.Options
	for i, opt := range settings.invokes {
		if opt == nil {
			settings.invokes[i] = fx.Options()
		}
	}

	app := fx.New(
		fx.Options(ctors...),
		fx.Options(settings.invokes...),

		fx.NopLogger,
	)

	// TODO: we probably should have a 'firewall' for Closing signal
	//  on this context, and implement closing logic through lifecycles
	//  correctly
	if err := app.Start(ctx); err != nil {
		// comment fx.NopLogger few lines above for easier debugging
		return nil, xerrors.Errorf("starting node: %w", err)
	}

	return app.Stop, nil
}

//
func NetworkName() (dtypes.NetworkName, error) {
	return "", nil
}

func newDbProc(cfg *config.DBConfig) (db.DBProcessInterface, error) {
	if cfg != nil && cfg.Conn != "" {
		return db.NewDbProcess(cfg)
	}
	return nil, xerrors.New("mysql conn do not set")
}
