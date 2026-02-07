package capability_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sergekukharev/agent-samwise/internal/capability"
	"github.com/sergekukharev/agent-samwise/internal/config"
	"github.com/sergekukharev/agent-samwise/internal/output"
	"github.com/sergekukharev/agent-samwise/internal/platform"
)

type stubCalendarReader struct {
	events []platform.CalendarEvent
	err    error
}

func (s *stubCalendarReader) TodayEvents(calendarID string) ([]platform.CalendarEvent, error) {
	return s.events, s.err
}

type stubTaskCreator struct {
	created []platform.TodoistTask
	err     error
}

func (s *stubTaskCreator) CreateTask(task platform.TodoistTask) error {
	if s.err != nil {
		return s.err
	}
	s.created = append(s.created, task)
	return nil
}

func testConfig() config.Config {
	return config.Config{
		Calendar: config.CalendarConfig{CalendarID: "test-calendar"},
		Todoist:  config.TodoistConfig{ProjectID: "test-project"},
	}
}

func TestCalendarSync_NoEvents(t *testing.T) {
	var buf bytes.Buffer
	cs := &capability.CalendarSync{
		Calendar: &stubCalendarReader{},
		Todoist:  &stubTaskCreator{},
	}

	err := cs.Run(testConfig(), config.Secrets{}, &output.TerminalPresenter{Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No events today") {
		t.Errorf("output = %q, want 'No events today'", buf.String())
	}
}

func TestCalendarSync_CreatesTasksForEvents(t *testing.T) {
	todoist := &stubTaskCreator{}
	var buf bytes.Buffer

	startTime := time.Date(2026, 2, 6, 10, 0, 0, 0, time.UTC)
	cs := &capability.CalendarSync{
		Calendar: &stubCalendarReader{
			events: []platform.CalendarEvent{
				{Title: "Standup", StartTime: startTime, RSVP: platform.RSVPAccepted, MeetingLink: "https://meet.google.com/abc"},
				{Title: "Company Holiday", AllDay: true, RSVP: platform.RSVPAccepted},
			},
		},
		Todoist: todoist,
	}

	err := cs.Run(testConfig(), config.Secrets{}, &output.TerminalPresenter{Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(todoist.created) != 2 {
		t.Fatalf("created %d tasks, want 2", len(todoist.created))
	}

	// Timed event
	task := todoist.created[0]
	if task.Title != "Standup" {
		t.Errorf("title = %q, want %q", task.Title, "Standup")
	}
	if task.DueDateTime == nil {
		t.Fatal("DueDateTime should be set for timed event")
	}
	if task.DueDateTime.Hour() != 10 {
		t.Errorf("due hour = %d, want 10", task.DueDateTime.Hour())
	}
	if task.Description != "https://meet.google.com/abc" {
		t.Errorf("description = %q, want meeting link", task.Description)
	}
	if task.Priority != 3 {
		t.Errorf("priority = %d, want 3", task.Priority)
	}
	if task.ProjectID != "test-project" {
		t.Errorf("project_id = %q, want %q", task.ProjectID, "test-project")
	}

	// All-day event
	allDay := todoist.created[1]
	if allDay.Title != "Company Holiday" {
		t.Errorf("title = %q, want %q", allDay.Title, "Company Holiday")
	}
	if allDay.DueDateTime != nil {
		t.Error("DueDateTime should be nil for all-day event")
	}
}

func TestCalendarSync_FiltersDeclinedEvents(t *testing.T) {
	todoist := &stubTaskCreator{}
	var buf bytes.Buffer

	cs := &capability.CalendarSync{
		Calendar: &stubCalendarReader{
			events: []platform.CalendarEvent{
				{Title: "Accepted meeting", RSVP: platform.RSVPAccepted, AllDay: true},
				{Title: "Declined meeting", RSVP: platform.RSVPDeclined, AllDay: true},
				{Title: "Tentative meeting", RSVP: platform.RSVPTentative, AllDay: true},
			},
		},
		Todoist: todoist,
	}

	err := cs.Run(testConfig(), config.Secrets{}, &output.TerminalPresenter{Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(todoist.created) != 2 {
		t.Fatalf("created %d tasks, want 2 (declined should be filtered)", len(todoist.created))
	}
}

func TestCalendarSync_UnconfirmedPrefix(t *testing.T) {
	todoist := &stubTaskCreator{}
	var buf bytes.Buffer

	cs := &capability.CalendarSync{
		Calendar: &stubCalendarReader{
			events: []platform.CalendarEvent{
				{Title: "Design Review", RSVP: platform.RSVPNeedsAction, AllDay: true},
			},
		},
		Todoist: todoist,
	}

	err := cs.Run(testConfig(), config.Secrets{}, &output.TerminalPresenter{Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if todoist.created[0].Title != "UNCONFIRMED: Design Review" {
		t.Errorf("title = %q, want %q", todoist.created[0].Title, "UNCONFIRMED: Design Review")
	}
}

func TestCalendarSync_CalendarAPIError(t *testing.T) {
	var buf bytes.Buffer
	cs := &capability.CalendarSync{
		Calendar: &stubCalendarReader{err: fmt.Errorf("API unavailable")},
		Todoist:  &stubTaskCreator{},
	}

	err := cs.Run(testConfig(), config.Secrets{}, &output.TerminalPresenter{Writer: &buf})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "API unavailable") {
		t.Errorf("error = %q, want mention of API unavailable", err.Error())
	}
}

func TestCalendarSync_TodoistAPIError(t *testing.T) {
	var buf bytes.Buffer
	cs := &capability.CalendarSync{
		Calendar: &stubCalendarReader{
			events: []platform.CalendarEvent{
				{Title: "Meeting", RSVP: platform.RSVPAccepted, AllDay: true},
			},
		},
		Todoist: &stubTaskCreator{err: fmt.Errorf("rate limited")},
	}

	err := cs.Run(testConfig(), config.Secrets{}, &output.TerminalPresenter{Writer: &buf})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "rate limited") {
		t.Errorf("error = %q, want mention of rate limited", err.Error())
	}
}

func TestCalendarSync_OutputSummary(t *testing.T) {
	todoist := &stubTaskCreator{}
	var buf bytes.Buffer

	cs := &capability.CalendarSync{
		Calendar: &stubCalendarReader{
			events: []platform.CalendarEvent{
				{Title: "Standup", RSVP: platform.RSVPAccepted, AllDay: true},
				{Title: "Retro", RSVP: platform.RSVPAccepted, AllDay: true},
			},
		},
		Todoist: todoist,
	}

	err := cs.Run(testConfig(), config.Secrets{}, &output.TerminalPresenter{Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Created 2 tasks") {
		t.Errorf("output missing task count, got:\n%s", got)
	}
	if !strings.Contains(got, "- Standup") {
		t.Errorf("output missing task name, got:\n%s", got)
	}
}
