package main

import (
	"context"

	"github.com/urfave/cli/v2"
	"go.opencensus.io/trace"

	"github.com/ipfs-force-community/venus-messager/build"
	lcli "github.com/ipfs-force-community/venus-messager/cli"
	"github.com/ipfs-force-community/venus-messager/cmd"
	"github.com/ipfs-force-community/venus-messager/lib/log"
	"github.com/ipfs-force-community/venus-messager/lib/tracing"
	"github.com/ipfs-force-community/venus-messager/node/repo"
)

func main() {
	build.RunningNodeType = build.NodeFull

	log.SetupLogLevels()

	local := []*cli.Command{
		cmd.DaemonCmd,
	}

	jaeger := tracing.SetupJaegerTracing("venus-message")
	defer func() {
		if jaeger != nil {
			jaeger.Flush()
		}
	}()

	for _, cmd := range local {
		cmd := cmd
		originBefore := cmd.Before
		cmd.Before = func(cctx *cli.Context) error {
			trace.UnregisterExporter(jaeger)
			jaeger = tracing.SetupJaegerTracing("venus-message/" + cmd.Name)

			if originBefore != nil {
				return originBefore(cctx)
			}
			return nil
		}
	}
	ctx, span := trace.StartSpan(context.Background(), "/cli")
	defer span.End()

	app := &cli.App{
		Name:                 "venus-message",
		Usage:                "Filecoin decentralized storage network client",
		Version:              build.UserVersion(),
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "repo",
				EnvVars: []string{"VENUS_MESSAGE_PATH"},
				Hidden:  true,
				Value:   "~/.venus_message",
			},
		},

		Commands: local,
	}
	app.Setup()
	app.Metadata["traceContext"] = ctx
	app.Metadata["repoType"] = repo.FullNode

	lcli.RunApp(app)
}
