package main

import (
	"os"

	"github.com/sergekukharev/agent-samwise/internal/cli"
)

func main() {
	router := cli.NewRouter(capabilities())
	os.Exit(router.Run(os.Args[1:]))
}

func capabilities() []cli.Capability {
	return []cli.Capability{
		// Capabilities will be registered here as they are built.
	}
}
