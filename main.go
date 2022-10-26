package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	cli "github.com/jawher/mow.cli"

	cmd_Cli "github.com/jeffreylean/blaster/cmd/cli"
	cmd_Server "github.com/jeffreylean/blaster/cmd/server"
)

func main() {
	app := cli.App("blaster", "")

	app.Command("cli", "run in cli mode", cmd_Cli.CmdCLI)
	app.Command("server", "run in http server mode", cmd_Server.CmdServer)
	// Invoke the app passing in os.Args.
	app.Run(os.Args)
}

// Log logs the error
func log(action string, err error) {
	fmt.Printf("unable to %s due to %s\n", action, err.Error())
}
