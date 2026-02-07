package capability

import (
	"fmt"
	"strings"

	"github.com/sergekukharev/agent-samwise/internal/config"
	"github.com/sergekukharev/agent-samwise/internal/output"
	"github.com/sergekukharev/agent-samwise/internal/platform"
)

// CalendarSync fetches today's calendar events and creates Todoist tasks.
type CalendarSync struct {
	Calendar platform.CalendarReader
	Todoist  platform.TaskCreator
}

func (cs *CalendarSync) Run(cfg config.Config, secrets config.Secrets, out output.Presenter) error {
	events, err := cs.Calendar.TodayEvents(cfg.Calendar.CalendarID)
	if err != nil {
		return fmt.Errorf("fetching calendar events: %w", err)
	}

	actionable := filterDeclined(events)

	if len(actionable) == 0 {
		return out.Present(output.Briefing{
			Title:    "Calendar Sync",
			Sections: []output.Section{{Heading: "Result", Body: "No events today"}},
		})
	}

	var created []string
	for _, event := range actionable {
		task := toTodoistTask(event, cfg.Todoist.ProjectID)
		if err := cs.Todoist.CreateTask(task); err != nil {
			return fmt.Errorf("creating todoist task %q: %w", task.Title, err)
		}
		created = append(created, task.Title)
	}

	return out.Present(output.Briefing{
		Title: "Calendar Sync",
		Sections: []output.Section{
			{
				Heading: "Result",
				Body:    fmt.Sprintf("Created %d tasks", len(created)),
			},
			{
				Heading: "Tasks",
				Body:    formatTaskList(created),
			},
		},
	})
}

func filterDeclined(events []platform.CalendarEvent) []platform.CalendarEvent {
	var result []platform.CalendarEvent
	for _, e := range events {
		if e.RSVP != platform.RSVPDeclined {
			result = append(result, e)
		}
	}
	return result
}

func toTodoistTask(event platform.CalendarEvent, projectID string) platform.TodoistTask {
	title := event.Title
	if event.RSVP == platform.RSVPNeedsAction {
		title = "UNCONFIRMED: " + title
	}

	task := platform.TodoistTask{
		Title:       title,
		Description: event.MeetingLink,
		ProjectID:   projectID,
		Priority:    3,
	}

	if !event.AllDay && !event.StartTime.IsZero() {
		t := event.StartTime
		task.DueDateTime = &t
	}

	return task
}

func formatTaskList(titles []string) string {
	var lines []string
	for _, t := range titles {
		lines = append(lines, "- "+t)
	}
	return strings.Join(lines, "\n")
}
