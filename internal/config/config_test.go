package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sergekukharev/agent-samwise/internal/config"
)

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTestConfig(t, `
calendar:
  calendar_id: "primary"
todoist:
  project_id: "12345"
  kanban_board_id: "67890"
areas:
  - name: "Backend"
    keywords: ["go", "api", "backend"]
  - name: "Frontend"
    keywords: ["react", "ui"]
`)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Calendar.CalendarID != "primary" {
		t.Errorf("calendar_id = %q, want %q", cfg.Calendar.CalendarID, "primary")
	}
	if cfg.Todoist.ProjectID != "12345" {
		t.Errorf("project_id = %q, want %q", cfg.Todoist.ProjectID, "12345")
	}
	if cfg.Todoist.KanbanBoardID != "67890" {
		t.Errorf("kanban_board_id = %q, want %q", cfg.Todoist.KanbanBoardID, "67890")
	}
	if len(cfg.Areas) != 2 {
		t.Errorf("areas count = %d, want 2", len(cfg.Areas))
	}
	if cfg.Areas[0].Name != "Backend" {
		t.Errorf("areas[0].name = %q, want %q", cfg.Areas[0].Name, "Backend")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if got := err.Error(); got != "config file not found: /nonexistent/config.yaml" {
		t.Errorf("error = %q, want mention of file path", got)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTestConfig(t, `
calendar:
  calendar_id: [invalid
`)

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestValidateFor_Calendar_MissingID(t *testing.T) {
	cfg := config.Config{}
	err := cfg.ValidateFor("calendar")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateFor_Calendar_Valid(t *testing.T) {
	cfg := config.Config{
		Calendar: config.CalendarConfig{CalendarID: "primary"},
	}
	if err := cfg.ValidateFor("calendar"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateFor_UnrelatedCapability(t *testing.T) {
	// An empty config should pass validation for a capability with no requirements
	cfg := config.Config{}
	if err := cfg.ValidateFor("gmail"); err != nil {
		t.Errorf("unexpected error for gmail with empty config: %v", err)
	}
}

func TestResolveSecrets_AllPresent(t *testing.T) {
	t.Setenv("GOOGLE_CREDENTIALS", "google-creds")
	t.Setenv("TODOIST_API_TOKEN", "todoist-token")

	secrets, err := config.ResolveSecrets("calendar", "todoist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets.GoogleCredentials != "google-creds" {
		t.Errorf("GoogleCredentials = %q, want %q", secrets.GoogleCredentials, "google-creds")
	}
	if secrets.TodoistAPIToken != "todoist-token" {
		t.Errorf("TodoistAPIToken = %q, want %q", secrets.TodoistAPIToken, "todoist-token")
	}
}

func TestResolveSecrets_MissingVar(t *testing.T) {
	t.Setenv("GOOGLE_CREDENTIALS", "")

	_, err := config.ResolveSecrets("calendar")
	if err == nil {
		t.Fatal("expected error for missing env var")
	}
}

func TestResolveSecrets_OnlyChecksNeededCapabilities(t *testing.T) {
	// Don't set any env vars — but request no capabilities that need them
	secrets, err := config.ResolveSecrets()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets.GoogleCredentials != "" {
		t.Errorf("expected empty GoogleCredentials, got %q", secrets.GoogleCredentials)
	}
}

func writeTestConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing test config: %v", err)
	}
	return path
}
