package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sergekukharev/agent-samwise/internal/cli"
	"github.com/sergekukharev/agent-samwise/internal/config"
	"github.com/sergekukharev/agent-samwise/internal/output"
)

func main() {
	router := cli.NewRouter(capabilities())
	os.Exit(router.Run(os.Args[1:]))
}

func capabilities() []cli.Capability {
	return []cli.Capability{
		hello(),
	}
}

// hello is a placeholder capability that proves the full pipeline works.
// Remove once a real capability is implemented.
func hello() cli.Capability {
	return cli.Capability{
		Name:        "hello",
		Description: "Test the pipeline (temporary)",
		Run: func(cfg config.Config, secrets config.Secrets, out output.Presenter) error {
			return out.Present(output.Briefing{
				Title: "Good morning!",
				Sections: []output.Section{
					{
						Heading: "Status",
						Body:    fmt.Sprintf("Sam is running. Time: %s", time.Now().Format("15:04")),
					},
					{
						Heading: "Pipeline",
						Body:    "Config loaded, secrets resolved, presenter detected. All systems go.",
					},
				},
			})
		},
	}
}
