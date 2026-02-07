## Context

This is the first real capability for Sam. It needs to integrate with two external APIs: Google Calendar (read) and Todoist (write). Both APIs are REST-based. The capability runs as the `calendar-sync` subcommand and outputs a summary via the presenter.

## Goals / Non-Goals

**Goals:**
- Fetch today's events from Google Calendar
- Create Todoist tasks with formatted titles and meeting links
- Establish API client patterns that future capabilities will reuse
- Testable via interfaces — no real API calls in tests

**Non-Goals:**
- Deduplication (creating duplicate tasks is acceptable for now)
- Syncing future events
- Bidirectional sync (Todoist → Calendar)
- OAuth flow — credentials are provided via env vars

## Decisions

### Decision: Google Calendar API via REST
Using the Google Calendar v3 REST API directly with `net/http`, not the official Go client library.

**Rationale:** The official `google.golang.org/api/calendar/v3` library pulls in a large dependency tree (google-api-go-client, oauth2, cloud metadata). We only need one endpoint: `events.list` with `timeMin`/`timeMax`. A thin HTTP client with JSON parsing keeps dependencies minimal, consistent with the project's approach.

**Alternatives considered:**
- `google.golang.org/api/calendar/v3` — heavyweight, adds ~30 transitive dependencies for one API call

### Decision: Todoist REST API v2
Using the Todoist REST API v2 directly with `net/http`.

**Rationale:** Same reasoning — we only need `POST /tasks`. No Go SDK needed for a single endpoint. API token goes in the `Authorization: Bearer` header.

### Decision: Google Service Account credentials
Using a Google Service Account JSON key (via `GOOGLE_CREDENTIALS` env var) for authentication. The service account must have read access to the target calendar (shared with the service account email).

**Rationale:** Service accounts don't require user-interactive OAuth consent flows, which is essential for headless GitHub Actions runs. The JSON key can be stored as a GitHub Secret.

**Flow:**
1. Parse service account JSON from `GOOGLE_CREDENTIALS`
2. Generate a signed JWT requesting `https://www.googleapis.com/auth/calendar.readonly` scope
3. Exchange JWT for an access token via Google's token endpoint
4. Use access token for Calendar API calls

### Decision: Interfaces for API clients
Define interfaces in `internal/platform/` for both Google Calendar and Todoist. Concrete implementations use `net/http`. Tests use stubs.

```go
// internal/platform/calendar.go
type RSVPStatus string

const (
    RSVPAccepted    RSVPStatus = "accepted"
    RSVPDeclined    RSVPStatus = "declined"
    RSVPNeedsAction RSVPStatus = "needsAction"
    RSVPTentative   RSVPStatus = "tentative"
)

type CalendarEvent struct {
    Title       string
    StartTime   time.Time
    AllDay      bool
    MeetingLink string
    RSVP        RSVPStatus
}

type CalendarReader interface {
    TodayEvents(calendarID string) ([]CalendarEvent, error)
}

// internal/platform/todoist.go
type TodoistTask struct {
    Title       string
    Description string
    ProjectID   string
    DueDateTime *time.Time  // nil for all-day events
    Priority    int         // Todoist priority (1-4)
}

type TaskCreator interface {
    CreateTask(task TodoistTask) error
}
```

**Rationale:** Keeps the capability logic testable without real API calls. The capability function accepts interfaces, not concrete clients.

### Decision: Project layout for the capability

```
internal/
├── platform/
│   ├── calendar.go      # CalendarEvent, CalendarReader interface
│   ├── gcalendar.go     # Google Calendar REST implementation
│   ├── todoist.go       # TodoistTask, TaskCreator interface
│   └── todoist_rest.go  # Todoist REST implementation
└── capability/
    └── calendarsync.go  # Capability logic: fetch events → create tasks → present
```

The `capability/` package contains the business logic that wires platform clients to the presenter. `cmd/sam/main.go` registers it.

## Risks / Trade-offs

- **No deduplication** → Running `calendar-sync` twice creates duplicate tasks. Acceptable for now — user runs it once each morning. Can add idempotency later using event IDs stored in task metadata.
- **Service account access** → The target calendar must be explicitly shared with the service account email. If the user forgets this step, the API returns an empty event list or 403.
- **No token caching** → Each run generates a new JWT and exchanges it for an access token. Acceptable for a tool that runs once a day.
- **Timezone** → Events are filtered using the system's local timezone. In GitHub Actions this defaults to UTC. The user may need to set `TZ` env var in the workflow.

## Open Questions

- What timezone should be used in GitHub Actions? Defaulting to system timezone, but the workflow may need `TZ=Europe/Berlin` (or wherever the user is) to get the right "today".
