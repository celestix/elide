package cmd

import (
	"fmt"
	"os"
	"strings"
)

const info = `Elide V1
A multi-protocol extensional API to provide support for some missing methods from the official Telegram Bot API like resolveUsername, getMessages etc.`

const help = `
Commands:
- help: shows this help section
- init: creates a preset config in the current directory
- run: starts the API server
`

func Execute() {
	args := os.Args[1:]
	if len(args) == 0 {
		// a hacky way to prevent more lines of codes
		args = []string{"help"}
	}
	switch arg := strings.ToLower(args[0]); arg {
	case "init":
		createConfig()
		fmt.Println(`Elide: Created a preset "elide.yaml", please edit it with your telegram auth keys!`)
	case "run":
		loadConfig()
		runServer()
	case "help":
		fmt.Printf("%s\n%s", info, help)
	default:
		fmt.Printf("'%s' is not a valid command!\n%s", arg, help)
	}
}
