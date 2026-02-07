package platform

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTodoistClient_CreateTask(t *testing.T) {
	var received todoistCreateTaskRequest
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &TodoistClient{
		apiToken:   "test-token",
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	dueTime := time.Date(2026, 2, 6, 10, 0, 0, 0, time.UTC)
	err := client.CreateTask(TodoistTask{
		Title:       "Standup",
		Description: "https://meet.google.com/abc",
		ProjectID:   "12345",
		DueDateTime: &dueTime,
		Priority:    3,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedAuth != "Bearer test-token" {
		t.Errorf("auth = %q, want %q", receivedAuth, "Bearer test-token")
	}
	if received.Content != "Standup" {
		t.Errorf("content = %q, want %q", received.Content, "Standup")
	}
	if received.Description != "https://meet.google.com/abc" {
		t.Errorf("description = %q, want meeting link", received.Description)
	}
	if received.ProjectID != "12345" {
		t.Errorf("project_id = %q, want %q", received.ProjectID, "12345")
	}
	if received.Priority != 3 {
		t.Errorf("priority = %d, want 3", received.Priority)
	}
	if received.DueDatetime == nil {
		t.Fatal("due_datetime should be set")
	}
}

func TestTodoistClient_CreateTask_NoDueTime(t *testing.T) {
	var received todoistCreateTaskRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &TodoistClient{
		apiToken:   "test-token",
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	err := client.CreateTask(TodoistTask{
		Title:     "Company Holiday",
		ProjectID: "12345",
		Priority:  3,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.DueDatetime != nil {
		t.Error("due_datetime should be nil for all-day events")
	}
}

func TestTodoistClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden"}`))
	}))
	defer server.Close()

	client := &TodoistClient{
		apiToken:   "bad-token",
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	err := client.CreateTask(TodoistTask{
		Title:     "Test",
		ProjectID: "12345",
		Priority:  3,
	})

	if err == nil {
		t.Fatal("expected error for 403 response")
	}
}
