package main

import (
	_ "net/http/pprof"
	"os"

	cli "github.com/jawher/mow.cli"
	blastCli "github.com/jeffreylean/blaster/cmd/cli"
	blastServer "github.com/jeffreylean/blaster/cmd/server"
)

func main() {
	app := cli.App("blaster", "")

	app.Command("cli", "run in cli mode", blastCli.CmdCLI)
	app.Command("server", "run in http server mode", blastServer.CmdServer)
	// Invoke the app passing in os.Args.
	app.Run(os.Args)
}
