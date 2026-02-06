package cli

import (
	"github.com/sergekukharev/agent-samwise/internal/config"
	"github.com/sergekukharev/agent-samwise/internal/output"
)

// Capability represents a Sam skill that can be invoked as a subcommand.
type Capability struct {
	Name        string
	Description string
	Run         func(cfg config.Config, secrets config.Secrets, out output.Presenter) error
	// RequiredEnv lists the capability names used to resolve secrets
	// (e.g., "calendar", "todoist", "slack").
	RequiredEnv []string
	// RequiredConfig lists the capability names used to validate config sections.
	RequiredConfig []string
}
