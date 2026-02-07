## Why

Every morning I want to see today's meetings in my Todoist task list alongside my other tasks. This gives me a unified view of my day — what I need to do and when I'll be in meetings. Currently I have to check my calendar separately, which means I either forget about meetings or can't plan work around them.

## What Changes

- New `calendar-sync` subcommand that fetches today's events from a configured Google Calendar and creates a Todoist task for each one
- Google Calendar API client to list today's events
- Todoist API client to create tasks in a configured project
- Task title includes event time and name (e.g., "10:00 - Standup")
- Task description includes meeting link if present
- No due dates — tasks are just a visual reference in the project
- No deduplication — always creates fresh tasks (can be improved later)

## Capabilities

### New Capabilities
- `calendar-sync`: Fetching today's Google Calendar events and creating corresponding Todoist tasks

### Modified Capabilities

## Impact

- New Google Calendar API dependency (REST API, OAuth credentials via env var)
- New Todoist API dependency (REST API, API token via env var)
- Registers `calendar-sync` subcommand in CLI
- Requires `GOOGLE_CREDENTIALS` and `TODOIST_API_TOKEN` environment variables
- Requires `calendar.calendar_id` and `todoist.project_id` in config
