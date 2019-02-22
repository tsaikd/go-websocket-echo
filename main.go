package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/tsaikd/KDGoLib/cliutil/cobrather"
	"github.com/tsaikd/go-websocket-echo/cmd/ping"
	"github.com/tsaikd/go-websocket-echo/cmd/server"
)

// Module info of package
var Module = &cobrather.Module{
	Use:   "go-websocket-echo",
	Short: "websocket echo client/server in golang",
	Commands: []*cobrather.Module{
		ping.Module,
		server.Module,
		cobrather.VersionModule,
	},
	RunE: func(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
		if len(args) < 1 {
			return cmd.Help()
		}
		return
	},
}

func main() {
	Module.MustMainRun()
}
