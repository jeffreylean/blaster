package main

import (
	"fmt"
	"os"
	"strconv"

	_ "net/http/pprof"

	cli "github.com/jawher/mow.cli"
	"github.com/jeffreylean/blaster/internal/blast"
	"github.com/jeffreylean/blaster/internal/server"
)

func main() {
	app := cli.App("blaster", "")

	app.Command("cli", "run in cli mode", cmdCLI)
	app.Command("server", "run in http server mode", cmdServer)
	// Invoke the app passing in os.Args
	app.Run(os.Args)
}

func cmdCLI(cmd *cli.Cmd) {
	cmd.Spec = "URI [-w] [-r] [-u] [--payload]"
	var (
		uri      = cmd.StringArg("URI", "", "The target server that you want to blast.")
		workers  = cmd.StringOpt("w workers", "", "The number of workers to work on your request.")
		requests = cmd.StringOpt("r requests", "", "The number of request to send.")
		rampup   = cmd.StringOpt("u ramp-up", "", "The duration for blaster to take to ramp-up to the full number of workers chosen.")
		payload  = cmd.StringOpt("payload", "", "HTTP body")
	)

	// Specify the action to execute when the app is invoked correctly
	cmd.Action = func() {
		// Build arguments
		w, err := strconv.Atoi(*workers)
		if err != nil {
			fmt.Println("Errors: ", err)
			os.Exit(2)
		}

		r, err := strconv.Atoi(*requests)
		if err != nil {
			fmt.Println("Errors: ", err)
			os.Exit(2)
		}

		u, err := strconv.Atoi(*rampup)
		if err != nil {
			fmt.Println("Errors: ", err)
			os.Exit(2)
		}

		// Send the command
		blast.Blast(*(uri), *(payload), int64(w), int64(r), int64(u))
	}
}

func cmdServer(cmd *cli.Cmd) {
	cmd.Spec = "PORT"

	var (
		port = cmd.StringArg("PORT", "", "Port the server should run on.")
	)

	// Specify the action to execute when the app is invoked correctly
	cmd.Action = func() {
		// Send the command
		server.Start(*port)
	}
}

// ---------------------------------------------------------------------------

// Log logs the error
func log(action string, err error) {
	fmt.Printf("unable to %s due to %s\n", action, err.Error())
}
