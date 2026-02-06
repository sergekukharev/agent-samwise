package output

import (
	"fmt"
	"os"
)

// DetectPresenter returns the appropriate Presenter based on the execution context.
// In GitHub Actions, it requires SLACK_WEBHOOK_URL and returns a SlackPresenter.
// Locally, it returns a TerminalPresenter.
func DetectPresenter() (Presenter, error) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		return NewTerminalPresenter(), nil
	}

	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return nil, fmt.Errorf("running in GitHub Actions but SLACK_WEBHOOK_URL is not set")
	}

	return NewSlackPresenter(webhookURL), nil
}
