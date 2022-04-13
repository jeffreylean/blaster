package main

import (
	"fmt"
	"os"
	"strconv"

	_ "net/http/pprof"

	cli "github.com/jawher/mow.cli"
	"github.com/jeffreylean/blaster/internal/blast"
)

func main() {
	app := cli.App("blaster", "")
	app.Spec = "URI [-w] [-r] [--payload]"
	var (
		uri      = app.StringArg("URI", "", "The target server that you want to blast.")
		workers  = app.StringOpt("w workers", "", "The number of workers to work on your request.")
		requests = app.StringOpt("r requests", "", "The number of request to send.")
		payload  = app.StringOpt("payload", "", "Json string")
	)

	// Specify the action to execute when the app is invoked correctly
	app.Action = func() {
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

		// Send the command
		blast.Blast(*(uri), *(payload), int64(w), int64(r))
	}
	// Invoke the app passing in os.Args
	app.Run(os.Args)
}

// ---------------------------------------------------------------------------

// Log logs the error
func log(action string, err error) {
	fmt.Printf("unable to %s due to %s\n", action, err.Error())
}
