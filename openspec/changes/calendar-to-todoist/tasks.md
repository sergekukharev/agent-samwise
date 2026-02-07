# Tasks

## 1. Platform Interfaces
- [x] 1.1 Define `CalendarEvent` value object (with `RSVPStatus`) and `CalendarReader` interface in `internal/platform/calendar.go`
- [x] 1.2 Define `TodoistTask` value object (with `DueDateTime` and `Priority`) and `TaskCreator` interface in `internal/platform/todoist.go`

## 2. Google Calendar Client
- [x] 2.1 Implement service account JWT signing and token exchange in `internal/platform/gcalendar.go`
- [x] 2.2 Implement `TodayEvents` — call Calendar v3 `events.list` with `timeMin`/`timeMax` for today
- [x] 2.3 Parse event response: extract title, start time, all-day flag, conference/meeting link, and RSVP status
- [x] 2.4 Write tests for event parsing (stub HTTP responses)

## 3. Todoist Client
- [x] 3.1 Implement `CreateTask` — POST to Todoist REST API v2 `/tasks` with Bearer token auth
- [x] 3.2 Write tests for task creation (stub HTTP responses)

## 4. Capability Logic
- [x] 4.1 Implement `calendarsync` capability in `internal/capability/calendarsync.go`: fetch events → filter declined → format titles (UNCONFIRMED prefix for needsAction) → create tasks with priority 3 and due time → present summary
- [x] 4.2 Handle edge cases: no events, all-day events, events with/without meeting links, declined events, unconfirmed events
- [x] 4.3 Write tests for capability logic using stub CalendarReader and TaskCreator

## 5. CLI Wiring
- [x] 5.1 Register `calendar-sync` subcommand in `cmd/sam/main.go` with required config and env vars
- [x] 5.2 Wire Google Calendar client and Todoist client to the capability
- [x] 5.3 Update GitHub Actions workflow to set `TZ` env var

## 6. Remove Hello Capability
- [x] 6.1 Remove the `hello` placeholder capability from `cmd/sam/main.go`
