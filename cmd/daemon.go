package cmd

import (
	"context"

	"github.com/mitchellh/go-homedir"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"

	"github.com/ipfs-force-community/venus-messager/api"
	"github.com/ipfs-force-community/venus-messager/build"
	"github.com/ipfs-force-community/venus-messager/metrics"
	"github.com/ipfs-force-community/venus-messager/node"
	"github.com/ipfs-force-community/venus-messager/node/config"
	"github.com/ipfs-force-community/venus-messager/node/modules/dtypes"
	"github.com/ipfs-force-community/venus-messager/node/repo"
)

var DaemonCmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start a daemon process",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "repo", Value: "~/.venus-message"},
		&cli.StringFlag{Name: "api", Value: "7979"},
		&cli.BoolFlag{Name: "bootstrap", Value: true},
		&cli.StringFlag{Name: "network", Value: ""},
		&cli.StringFlag{Name: "message-store", Value: "messagestore"},
		&cli.BoolFlag{Name: "debug-model", Value: true},
	},
	Action: func(cctx *cli.Context) error {
		ctx, _ := tag.New(context.Background(), tag.Insert(metrics.Version, build.BuildVersion))
		dir, err := homedir.Expand(cctx.String("repo"))
		if err != nil {
			log.Warnw("could not expand repo location", "error", err)
		} else {
			log.Infof("venus-messager repo: %s", dir)
		}
		r, err := repo.NewFS(cctx.String("repo"))
		if err != nil {
			return xerrors.Errorf("opening fs repo: %w", err)
		}

		if err := r.Init(repo.FullNode); err != nil && err != repo.ErrRepoExists {
			return xerrors.Errorf("repo init error: %w", err)
		}

		var api api.FullNode

		stop, err := node.New(ctx,
			node.FullAPI(&api),

			node.ApplyIf(func(s *node.Settings) bool { return !cctx.Bool("bootstrap") }), //node.Local(),

			node.Online(),
			node.Repo(r),

			node.ApplyIf(func(s *node.Settings) bool { return cctx.IsSet("api") },
				node.Override(node.SetApiEndpointKey, func(lr repo.LockedRepo) error {
					apima, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/" +
						cctx.String("api"))
					if err != nil {
						return err
					}
					return lr.SetAPIEndpoint(apima)
				})),
			node.Override(new(*config.DBConfig),
				&config.DBConfig{Conn: cctx.String("message-store"), Type: "sqlite", DebugMode: cctx.Bool("debug-model")}),
			node.ApplyIf(func(s *node.Settings) bool { return cctx.IsSet("network") }),
			node.Override(new(dtypes.NetworkName), dtypes.NetworkName(cctx.String("network"))),
		)
		if err != nil {
			return xerrors.Errorf("initializing node: %w", err)
		}

		// Register all metric views
		if err = view.Register(
			metrics.DefaultViews...,
		); err != nil {
			log.Fatalf("Cannot register the view: %v", err)
		}

		// Set the metric to one so it is published to the exporter
		stats.Record(ctx, metrics.LotusInfo.M(1))

		endpoint, err := r.APIEndpoint()
		if err != nil {
			return xerrors.Errorf("getting api endpoint: %w", err)
		}

		// TODO: properly parse api endpoint (or make it a URL)
		return ServeRPC(api, stop, endpoint)
	},
}
