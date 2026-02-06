package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sergekukharev/agent-samwise/internal/cli"
	"github.com/sergekukharev/agent-samwise/internal/config"
	"github.com/sergekukharev/agent-samwise/internal/output"
)

func testCapabilities() []cli.Capability {
	return []cli.Capability{
		{
			Name:        "greet",
			Description: "Say hello",
			Run: func(cfg config.Config, secrets config.Secrets, out output.Presenter) error {
				return out.Present(output.Briefing{
					Title:    "Hello",
					Sections: []output.Section{{Heading: "Greeting", Body: "Hi there!"}},
				})
			},
		},
		{
			Name:           "calendar-sync",
			Description:    "Sync calendar to todoist",
			RequiredConfig: []string{"calendar"},
			RequiredEnv:    []string{"calendar"},
			Run: func(cfg config.Config, secrets config.Secrets, out output.Presenter) error {
				return nil
			},
		},
	}
}

func TestRouter_NoArgs_PrintsHelp(t *testing.T) {
	router := cli.NewRouter(testCapabilities())
	code := router.Run([]string{})
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
}

func TestRouter_UnknownCommand(t *testing.T) {
	router := cli.NewRouter(testCapabilities())
	code := router.Run([]string{"nonexistent"})
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
}

func TestRouter_KnownCommand(t *testing.T) {
	cfgPath := writeMinimalConfig(t)

	router := cli.NewRouter(testCapabilities())
	code := router.Run([]string{"--config", cfgPath, "greet"})
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
}

func TestRouter_ConfigFlag(t *testing.T) {
	cfgPath := writeMinimalConfig(t)

	var ran bool
	router := cli.NewRouter([]cli.Capability{
		{
			Name:        "test",
			Description: "test command",
			Run: func(cfg config.Config, secrets config.Secrets, out output.Presenter) error {
				ran = true
				return nil
			},
		},
	})

	code := router.Run([]string{"--config", cfgPath, "test"})
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !ran {
		t.Error("capability did not run")
	}
}

func TestRouter_MissingConfig(t *testing.T) {
	router := cli.NewRouter(testCapabilities())
	code := router.Run([]string{"--config", "/nonexistent/config.yaml", "greet"})
	if code != 1 {
		t.Errorf("exit code = %d, want 1 for missing config", code)
	}
}

func TestRouter_ConfigValidationFailure(t *testing.T) {
	// calendar-sync requires calendar.calendar_id but our config has it empty
	cfgPath := writeMinimalConfig(t)

	router := cli.NewRouter(testCapabilities())
	code := router.Run([]string{"--config", cfgPath, "calendar-sync"})
	if code != 1 {
		t.Errorf("exit code = %d, want 1 for config validation failure", code)
	}
}

func TestRouter_Suggestions(t *testing.T) {
	router := cli.NewRouter(testCapabilities())
	// "gret" is close to "greet" — should suggest it
	code := router.Run([]string{"gret"})
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
}

func writeMinimalConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("# minimal config\n"), 0644); err != nil {
		t.Fatalf("writing test config: %v", err)
	}
	return path
}
