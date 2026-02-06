package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultPath = "config.yaml"

// Config is the top-level configuration for Sam.
type Config struct {
	Calendar CalendarConfig `yaml:"calendar"`
	Todoist  TodoistConfig  `yaml:"todoist"`
	Gmail    GmailConfig    `yaml:"gmail"`
	Slack    SlackConfig    `yaml:"slack"`
	Areas    []Area         `yaml:"areas"`
}

type CalendarConfig struct {
	CalendarID string `yaml:"calendar_id"`
}

type TodoistConfig struct {
	ProjectID    string `yaml:"project_id"`
	KanbanBoardID string `yaml:"kanban_board_id"`
}

type GmailConfig struct {
	// No config fields yet — Gmail capability will use the authenticated user's inbox.
}

type SlackConfig struct {
	// Webhook URL comes from SLACK_WEBHOOK_URL env var, not config.
}

// Area represents a project or area of interest for calendar recommendations.
type Area struct {
	Name     string   `yaml:"name"`
	Keywords []string `yaml:"keywords"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config file not found: %s", path)
		}
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config file: %w", err)
	}

	return cfg, nil
}

// Secrets holds API keys and tokens resolved from environment variables.
type Secrets struct {
	GoogleCredentials string
	TodoistAPIToken   string
	SlackWebhookURL   string
}

// EnvVar defines a required environment variable for a capability.
type EnvVar struct {
	Name       string
	Capability string
}

var googleEnv = []EnvVar{
	{Name: "GOOGLE_CREDENTIALS", Capability: "calendar/gmail"},
}

var todoistEnv = []EnvVar{
	{Name: "TODOIST_API_TOKEN", Capability: "todoist"},
}

var slackEnv = []EnvVar{
	{Name: "SLACK_WEBHOOK_URL", Capability: "slack"},
}

// ResolveSecrets reads required environment variables and returns Secrets.
// Only the env vars needed by the given capabilities are checked.
func ResolveSecrets(capabilities ...string) (Secrets, error) {
	needed := requiredEnvVars(capabilities)

	var missing []string
	values := make(map[string]string)

	for _, ev := range needed {
		val := os.Getenv(ev.Name)
		if val == "" {
			missing = append(missing, fmt.Sprintf("%s (required by %s)", ev.Name, ev.Capability))
		}
		values[ev.Name] = val
	}

	if len(missing) > 0 {
		return Secrets{}, fmt.Errorf("missing environment variables:\n  %s", strings.Join(missing, "\n  "))
	}

	return Secrets{
		GoogleCredentials: values["GOOGLE_CREDENTIALS"],
		TodoistAPIToken:   values["TODOIST_API_TOKEN"],
		SlackWebhookURL:   values["SLACK_WEBHOOK_URL"],
	}, nil
}

func requiredEnvVars(capabilities []string) []EnvVar {
	var result []EnvVar
	for _, cap := range capabilities {
		switch cap {
		case "calendar", "gmail":
			result = append(result, googleEnv...)
		case "todoist":
			result = append(result, todoistEnv...)
		case "slack":
			result = append(result, slackEnv...)
		}
	}
	return dedupEnvVars(result)
}

func dedupEnvVars(vars []EnvVar) []EnvVar {
	seen := make(map[string]bool)
	var result []EnvVar
	for _, v := range vars {
		if !seen[v.Name] {
			seen[v.Name] = true
			result = append(result, v)
		}
	}
	return result
}

// ValidateFor checks that the config has the required sections for the given capability.
func (c Config) ValidateFor(capability string) error {
	switch capability {
	case "calendar":
		if c.Calendar.CalendarID == "" {
			return fmt.Errorf("calendar.calendar_id is required for the calendar capability")
		}
	case "todoist":
		if c.Todoist.ProjectID == "" {
			return fmt.Errorf("todoist.project_id is required for the todoist capability")
		}
	case "review-projects":
		if c.Todoist.KanbanBoardID == "" {
			return fmt.Errorf("todoist.kanban_board_id is required for the review-projects capability")
		}
	case "calendar-recommendations":
		if len(c.Areas) == 0 {
			return fmt.Errorf("areas is required for the calendar-recommendations capability")
		}
	}
	return nil
}
