package cmd

import (
	"fmt"
	"os"
	"strings"
)

func Execute() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}
	switch strings.ToLower(args[0]) {
	case "init":
		createConfig()
		fmt.Println(`Elide: Created a preset "elide.yaml", please edit it with your telegram auth keys!`)
	case "run":
		loadConfig()
		runServer()
	}
}
