package main

import (
	"os"

	"github.com/sergekukharev/agent-samwise/internal/capability"
	"github.com/sergekukharev/agent-samwise/internal/cli"
	"github.com/sergekukharev/agent-samwise/internal/config"
	"github.com/sergekukharev/agent-samwise/internal/output"
	"github.com/sergekukharev/agent-samwise/internal/platform"
)

func main() {
	router := cli.NewRouter(capabilities())
	os.Exit(router.Run(os.Args[1:]))
}

func capabilities() []cli.Capability {
	return []cli.Capability{
		calendarSync(),
	}
}

func calendarSync() cli.Capability {
	return cli.Capability{
		Name:           "calendar-sync",
		Description:    "Sync today's calendar events to Todoist",
		RequiredConfig: []string{"calendar", "todoist"},
		RequiredEnv:    []string{"calendar", "todoist"},
		Run: func(cfg config.Config, secrets config.Secrets, out output.Presenter) error {
			calendarClient, err := platform.NewGoogleCalendarClient(secrets.GoogleCredentials)
			if err != nil {
				return err
			}

			todoistClient := platform.NewTodoistClient(secrets.TodoistAPIToken)

			cs := &capability.CalendarSync{
				Calendar: calendarClient,
				Todoist:  todoistClient,
			}

			return cs.Run(cfg, secrets, out)
		},
	}
}
