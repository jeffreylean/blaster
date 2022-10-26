package cli

import (
	"fmt"
	"os"
	"strconv"

	cli "github.com/jawher/mow.cli"
	"github.com/jeffreylean/blaster/internal/blast"
	"github.com/jeffreylean/blaster/pkg/futil"
)

func CmdCLI(cmd *cli.Cmd) {
	cmd.Spec = "URI [-w] [-r] [-u] [--payload]"
	var (
		uri      = cmd.StringArg("URI", "", "The target server that you want to blast.")
		workers  = cmd.StringOpt("w workers", "", "The number of workers to work on your request.")
		requests = cmd.StringOpt("r requests", "", "The number of request to send.")
		rampup   = cmd.StringOpt("u ramp-up(s)", "", "The duration in seconds for blaster to take to ramp-up to the full number of workers chosen.")
		payload  = cmd.StringOpt("payload", "", "File path which contain JSON payload of HTTP body.")
	)

	// Specify the action to execute when the app is invoked correctly.
	cmd.Action = func() {
		// Build arguments.
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

		// Get payload from file.
		var p []byte
		if *(payload) != "" {
			b, err := futil.ReadFile(*(payload))
			if err != nil {
				fmt.Println("Errors: ", err)
				os.Exit(2)
			}
			p = b
		}

		// Send the command
		blast.Blast(*(uri), p, int64(w), int64(r), int64(u))
	}
}
