package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackPresenter delivers briefings via a Slack Incoming Webhook.
type SlackPresenter struct {
	WebhookURL string
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

// NewSlackPresenter creates a presenter that posts to the given webhook URL.
func NewSlackPresenter(webhookURL string) *SlackPresenter {
	return &SlackPresenter{
		WebhookURL: webhookURL,
		HTTPClient: http.DefaultClient,
	}
}

func (p *SlackPresenter) Present(briefing Briefing) error {
	payload := buildSlackPayload(briefing)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling slack payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating slack request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending slack message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack webhook returned status %d", resp.StatusCode)
	}

	return nil
}

type slackPayload struct {
	Blocks []slackBlock `json:"blocks"`
}

type slackBlock struct {
	Type string     `json:"type"`
	Text *slackText `json:"text,omitempty"`
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func buildSlackPayload(briefing Briefing) slackPayload {
	// Build the full message as markdown, delivered in a single section block.
	// Slack's mrkdwn is close enough to standard markdown for headings, bold, lists.
	var md string
	md += fmt.Sprintf("*%s*\n\n", briefing.Title)

	for i, section := range briefing.Sections {
		md += fmt.Sprintf("*%s*\n%s\n", section.Heading, section.Body)
		if i < len(briefing.Sections)-1 {
			md += "\n"
		}
	}

	return slackPayload{
		Blocks: []slackBlock{
			{
				Type: "section",
				Text: &slackText{Type: "mrkdwn", Text: md},
			},
		},
	}
}
