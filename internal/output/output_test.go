package output_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sergekukharev/agent-samwise/internal/output"
)

var testBriefing = output.Briefing{
	Title: "Morning Briefing",
	Sections: []output.Section{
		{Heading: "Calendar", Body: "3 events today"},
		{Heading: "Tasks", Body: "2 items overdue\n1 due today"},
	},
}

func TestTerminalPresenter(t *testing.T) {
	var buf bytes.Buffer
	p := &output.TerminalPresenter{Writer: &buf}

	if err := p.Present(testBriefing); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()

	if !strings.Contains(got, "# Morning Briefing") {
		t.Error("missing markdown title")
	}
	if !strings.Contains(got, "## Calendar") {
		t.Error("missing Calendar section heading")
	}
	if !strings.Contains(got, "3 events today") {
		t.Error("missing Calendar body")
	}
	if !strings.Contains(got, "## Tasks") {
		t.Error("missing Tasks section heading")
	}
	if !strings.Contains(got, "2 items overdue") {
		t.Error("missing Tasks body")
	}
}

type stubHTTPClient struct {
	request    *http.Request
	body       []byte
	statusCode int
}

func (c *stubHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.request = req
	body, _ := io.ReadAll(req.Body)
	c.body = body
	return &http.Response{
		StatusCode: c.statusCode,
		Body:       io.NopCloser(strings.NewReader("ok")),
	}, nil
}

func TestSlackPresenter_PayloadFormat(t *testing.T) {
	client := &stubHTTPClient{statusCode: http.StatusOK}
	p := &output.SlackPresenter{
		WebhookURL: "https://hooks.slack.com/test",
		HTTPClient: client,
	}

	if err := p.Present(testBriefing); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.request.Header.Get("Content-Type") != "application/json" {
		t.Error("expected Content-Type application/json")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(client.body, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}

	blocks, ok := payload["blocks"].([]interface{})
	if !ok {
		t.Fatal("missing blocks array")
	}

	// Single section block with full markdown content
	if len(blocks) != 1 {
		t.Errorf("blocks count = %d, want 1", len(blocks))
	}

	section := blocks[0].(map[string]interface{})
	if section["type"] != "section" {
		t.Errorf("block type = %q, want section", section["type"])
	}

	text := section["text"].(map[string]interface{})
	content := text["text"].(string)
	if !strings.Contains(content, "*Morning Briefing*") {
		t.Error("missing markdown title in Slack payload")
	}
	if !strings.Contains(content, "*Calendar*") {
		t.Error("missing Calendar heading in Slack payload")
	}
	if !strings.Contains(content, "3 events today") {
		t.Error("missing Calendar body in Slack payload")
	}
}

func TestSlackPresenter_NonOKStatus(t *testing.T) {
	client := &stubHTTPClient{statusCode: http.StatusInternalServerError}
	p := &output.SlackPresenter{
		WebhookURL: "https://hooks.slack.com/test",
		HTTPClient: client,
	}

	err := p.Present(testBriefing)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention status code, got: %v", err)
	}
}

func TestDetectPresenter_Local(t *testing.T) {
	t.Setenv("GITHUB_ACTIONS", "")

	p, err := output.DetectPresenter()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := p.(*output.TerminalPresenter); !ok {
		t.Errorf("expected TerminalPresenter, got %T", p)
	}
}

func TestDetectPresenter_GitHubActionsWithWebhook(t *testing.T) {
	t.Setenv("GITHUB_ACTIONS", "true")
	t.Setenv("SLACK_WEBHOOK_URL", "https://hooks.slack.com/test")

	p, err := output.DetectPresenter()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := p.(*output.SlackPresenter); !ok {
		t.Errorf("expected SlackPresenter, got %T", p)
	}
}

func TestDetectPresenter_GitHubActionsWithoutWebhook(t *testing.T) {
	t.Setenv("GITHUB_ACTIONS", "true")
	t.Setenv("SLACK_WEBHOOK_URL", "")

	_, err := output.DetectPresenter()
	if err == nil {
		t.Fatal("expected error for missing webhook in GitHub Actions")
	}
	if !strings.Contains(err.Error(), "SLACK_WEBHOOK_URL") {
		t.Errorf("error should mention SLACK_WEBHOOK_URL, got: %v", err)
	}
}
