package platform

import "time"

// TodoistTask represents a task to be created in Todoist.
type TodoistTask struct {
	Title       string
	Description string
	ProjectID   string
	DueDateTime *time.Time // nil for all-day events or tasks without a specific time
	Priority    int        // Todoist priority: 1 (normal) to 4 (urgent)
}

// TaskCreator creates tasks in Todoist.
type TaskCreator interface {
	CreateTask(task TodoistTask) error
}
