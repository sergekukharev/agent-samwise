package platform

import "time"

// RSVPStatus represents the user's response status to a calendar event.
type RSVPStatus string

const (
	RSVPAccepted    RSVPStatus = "accepted"
	RSVPDeclined    RSVPStatus = "declined"
	RSVPNeedsAction RSVPStatus = "needsAction"
	RSVPTentative   RSVPStatus = "tentative"
)

// CalendarEvent represents a single event from a calendar.
type CalendarEvent struct {
	Title       string
	StartTime   time.Time
	AllDay      bool
	MeetingLink string
	RSVP        RSVPStatus
}

// CalendarReader fetches events from a calendar.
type CalendarReader interface {
	TodayEvents(calendarID string) ([]CalendarEvent, error)
}
