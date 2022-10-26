package server

import (
	cli "github.com/jawher/mow.cli"
	"github.com/jeffreylean/blaster/internal/server"
)

func CmdServer(cmd *cli.Cmd) {
	cmd.Spec = "PORT"

	var (
		port = cmd.StringArg("PORT", "", "Port the server should run on.")
	)

	// Specify the action to execute when the app is invoked correctly.
	cmd.Action = func() {
		// Send the command.
		server.Start(*port)
	}
}
