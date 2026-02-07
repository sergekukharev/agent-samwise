package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TodoistClient creates tasks via the Todoist REST API v2.
type TodoistClient struct {
	apiToken   string
	baseURL    string
	httpClient *http.Client
}

// NewTodoistClient creates a client with the given API token.
func NewTodoistClient(apiToken string) *TodoistClient {
	return &TodoistClient{
		apiToken:   apiToken,
		baseURL:    "https://api.todoist.com/rest/v2",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *TodoistClient) CreateTask(task TodoistTask) error {
	payload := todoistCreateTaskRequest{
		Content:     task.Title,
		Description: task.Description,
		ProjectID:   task.ProjectID,
		Priority:    task.Priority,
	}

	if task.DueDateTime != nil {
		s := task.DueDateTime.Format(time.RFC3339)
		payload.DueDatetime = &s
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling task: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/tasks", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating todoist request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("creating todoist task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Todoist API returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

type todoistCreateTaskRequest struct {
	Content     string  `json:"content"`
	Description string  `json:"description,omitempty"`
	ProjectID   string  `json:"project_id"`
	Priority    int     `json:"priority"`
	DueDatetime *string `json:"due_datetime,omitempty"`
}
