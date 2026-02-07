## ADDED Requirements

### Requirement: Fetch today's calendar events
The system SHALL fetch all events from the configured Google Calendar for today (midnight to midnight in the user's local timezone).

#### Scenario: Calendar has events today
- **WHEN** the configured calendar has events scheduled for today
- **THEN** all non-declined events are returned with their title, start time, RSVP status, and conference/meeting link

#### Scenario: Declined events are excluded
- **WHEN** the user has declined a calendar event
- **THEN** that event is not included in the results

#### Scenario: Calendar has no events today
- **WHEN** the configured calendar has no events scheduled for today
- **THEN** an empty list is returned
- **AND** the output reports that no events were found

#### Scenario: All-day events
- **WHEN** the calendar has all-day events
- **THEN** they are included with "All day" as the time instead of a specific time

#### Scenario: Multi-day events spanning today
- **WHEN** a multi-day event includes today
- **THEN** it is included in today's events

### Requirement: Create Todoist tasks from events
The system SHALL create one Todoist task per calendar event in the configured project with priority 3.

#### Scenario: Accepted event with a specific time
- **WHEN** an event has a start time of 10:00, title "Standup", and user has accepted
- **THEN** a Todoist task is created with title "Standup", due time set to 10:00, and priority 3

#### Scenario: Unconfirmed event
- **WHEN** an event has title "Design Review" and user has not responded (needsAction)
- **THEN** a Todoist task is created with title "UNCONFIRMED: Design Review" and priority 3

#### Scenario: All-day event
- **WHEN** an event is an all-day event with title "Company Holiday"
- **THEN** a Todoist task is created with title "Company Holiday", no due time, and priority 3

#### Scenario: Event with a meeting link
- **WHEN** an event has a conference link (Google Meet, Zoom, etc.)
- **THEN** the Todoist task description contains the meeting link

#### Scenario: Event without a meeting link
- **WHEN** an event has no conference link
- **THEN** the Todoist task is created with no description

### Requirement: Output summary
The system SHALL report what it did via the presenter.

#### Scenario: Events synced successfully
- **WHEN** calendar events are synced to Todoist
- **THEN** the output lists each created task with its title

#### Scenario: No events to sync
- **WHEN** there are no events today
- **THEN** the output reports "No events today"

### Requirement: Error handling
The system SHALL report clear errors when external APIs fail.

#### Scenario: Google Calendar API failure
- **WHEN** the Google Calendar API returns an error
- **THEN** the application exits with an error message describing the failure

#### Scenario: Todoist API failure
- **WHEN** the Todoist API returns an error while creating a task
- **THEN** the application exits with an error message describing the failure
- **AND** any tasks already created in this run remain (no rollback)
