package platform

import (
	"testing"
	"time"
)

func TestParseCalendarEvents_TimedEvent(t *testing.T) {
	items := []calendarEventItem{
		{
			Summary: "Standup",
			Start:   calendarEventTime{DateTime: "2026-02-06T10:00:00+01:00"},
			Attendees: []calendarAttendee{
				{Email: "me@example.com", Self: true, ResponseStatus: "accepted"},
			},
		},
	}

	events := parseCalendarEvents(items)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e := events[0]
	if e.Title != "Standup" {
		t.Errorf("title = %q, want %q", e.Title, "Standup")
	}
	if e.AllDay {
		t.Error("expected AllDay = false")
	}
	if e.StartTime.Hour() != 10 {
		t.Errorf("start hour = %d, want 10", e.StartTime.Hour())
	}
	if e.RSVP != RSVPAccepted {
		t.Errorf("RSVP = %q, want %q", e.RSVP, RSVPAccepted)
	}
}

func TestParseCalendarEvents_AllDayEvent(t *testing.T) {
	items := []calendarEventItem{
		{
			Summary: "Company Holiday",
			Start:   calendarEventTime{Date: "2026-02-06"},
		},
	}

	events := parseCalendarEvents(items)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e := events[0]
	if !e.AllDay {
		t.Error("expected AllDay = true")
	}
	if e.StartTime != (time.Time{}) {
		t.Errorf("expected zero StartTime for all-day event, got %v", e.StartTime)
	}
}

func TestParseCalendarEvents_DeclinedRSVP(t *testing.T) {
	items := []calendarEventItem{
		{
			Summary: "Meeting I declined",
			Start:   calendarEventTime{DateTime: "2026-02-06T14:00:00+01:00"},
			Attendees: []calendarAttendee{
				{Email: "me@example.com", Self: true, ResponseStatus: "declined"},
			},
		},
	}

	events := parseCalendarEvents(items)
	if events[0].RSVP != RSVPDeclined {
		t.Errorf("RSVP = %q, want %q", events[0].RSVP, RSVPDeclined)
	}
}

func TestParseCalendarEvents_NeedsActionRSVP(t *testing.T) {
	items := []calendarEventItem{
		{
			Summary: "Pending invite",
			Start:   calendarEventTime{DateTime: "2026-02-06T15:00:00+01:00"},
			Attendees: []calendarAttendee{
				{Email: "other@example.com", Self: false, ResponseStatus: "accepted"},
				{Email: "me@example.com", Self: true, ResponseStatus: "needsAction"},
			},
		},
	}

	events := parseCalendarEvents(items)
	if events[0].RSVP != RSVPNeedsAction {
		t.Errorf("RSVP = %q, want %q", events[0].RSVP, RSVPNeedsAction)
	}
}

func TestParseCalendarEvents_NoAttendees(t *testing.T) {
	items := []calendarEventItem{
		{
			Summary: "My own event",
			Start:   calendarEventTime{DateTime: "2026-02-06T09:00:00+01:00"},
		},
	}

	events := parseCalendarEvents(items)
	if events[0].RSVP != RSVPAccepted {
		t.Errorf("RSVP = %q, want %q for event with no attendees", events[0].RSVP, RSVPAccepted)
	}
}

func TestExtractMeetingLink_ConferenceData(t *testing.T) {
	item := calendarEventItem{
		ConferenceData: &conferenceData{
			EntryPoints: []conferenceEntryPoint{
				{EntryPointType: "video", URI: "https://meet.google.com/abc-defg-hij"},
			},
		},
	}

	link := extractMeetingLink(item)
	if link != "https://meet.google.com/abc-defg-hij" {
		t.Errorf("link = %q, want Google Meet URL", link)
	}
}

func TestExtractMeetingLink_HangoutFallback(t *testing.T) {
	item := calendarEventItem{
		HangoutLink: "https://meet.google.com/fallback",
	}

	link := extractMeetingLink(item)
	if link != "https://meet.google.com/fallback" {
		t.Errorf("link = %q, want hangout fallback URL", link)
	}
}

func TestExtractMeetingLink_NoLink(t *testing.T) {
	item := calendarEventItem{}

	link := extractMeetingLink(item)
	if link != "" {
		t.Errorf("link = %q, want empty string", link)
	}
}

func TestExtractRSVP_TentativeStatus(t *testing.T) {
	attendees := []calendarAttendee{
		{Email: "me@example.com", Self: true, ResponseStatus: "tentative"},
	}

	rsvp := extractRSVP(attendees)
	if rsvp != RSVPTentative {
		t.Errorf("RSVP = %q, want %q", rsvp, RSVPTentative)
	}
}
