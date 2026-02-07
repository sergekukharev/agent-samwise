package platform

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	calendarAPIBase   = "https://www.googleapis.com/calendar/v3"
	calendarReadScope = "https://www.googleapis.com/auth/calendar.readonly"
)

// ServiceAccountKey represents the JSON key file for a Google service account.
type ServiceAccountKey struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

// GoogleCalendarClient reads events from Google Calendar using a service account.
type GoogleCalendarClient struct {
	credentials ServiceAccountKey
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// NewGoogleCalendarClient creates a client from a service account JSON key.
func NewGoogleCalendarClient(credentialsJSON string) (*GoogleCalendarClient, error) {
	var key ServiceAccountKey
	if err := json.Unmarshal([]byte(credentialsJSON), &key); err != nil {
		return nil, fmt.Errorf("parsing service account credentials: %w", err)
	}

	if key.ClientEmail == "" || key.PrivateKey == "" {
		return nil, fmt.Errorf("service account credentials missing client_email or private_key")
	}

	if key.TokenURI == "" {
		key.TokenURI = googleTokenURL
	}

	return &GoogleCalendarClient{
		credentials: key,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *GoogleCalendarClient) TodayEvents(calendarID string) ([]CalendarEvent, error) {
	token, err := c.ensureToken()
	if err != nil {
		return nil, fmt.Errorf("authenticating with Google: %w", err)
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	eventsURL := fmt.Sprintf("%s/calendars/%s/events?%s",
		calendarAPIBase,
		url.PathEscape(calendarID),
		url.Values{
			"timeMin":      {startOfDay.Format(time.RFC3339)},
			"timeMax":      {endOfDay.Format(time.RFC3339)},
			"singleEvents": {"true"},
			"orderBy":      {"startTime"},
		}.Encode(),
	)

	req, err := http.NewRequest(http.MethodGet, eventsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating calendar request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching calendar events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google Calendar API returned %d: %s", resp.StatusCode, string(body))
	}

	var result calendarListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("parsing calendar response: %w", err)
	}

	return parseCalendarEvents(result.Items), nil
}

func (c *GoogleCalendarClient) ensureToken() (string, error) {
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   c.credentials.ClientEmail,
		"scope": calendarReadScope,
		"aud":   c.credentials.TokenURI,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}

	privateKey, err := parseRSAPrivateKey(c.credentials.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("parsing private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedJWT, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing JWT: %w", err)
	}

	resp, err := c.httpClient.PostForm(c.credentials.TokenURI, url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  {signedJWT},
	})
	if err != nil {
		return "", fmt.Errorf("exchanging JWT for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token exchange returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("parsing token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return c.accessToken, nil
}

func parseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	return rsaKey, nil
}

// Google Calendar API response types

type calendarListResponse struct {
	Items []calendarEventItem `json:"items"`
}

type calendarEventItem struct {
	Summary        string              `json:"summary"`
	Start          calendarEventTime   `json:"start"`
	End            calendarEventTime   `json:"end"`
	ConferenceData *conferenceData     `json:"conferenceData"`
	HangoutLink    string              `json:"hangoutLink"`
	Attendees      []calendarAttendee  `json:"attendees"`
}

type calendarEventTime struct {
	DateTime string `json:"dateTime"`
	Date     string `json:"date"`
}

type conferenceData struct {
	EntryPoints []conferenceEntryPoint `json:"entryPoints"`
}

type conferenceEntryPoint struct {
	EntryPointType string `json:"entryPointType"`
	URI            string `json:"uri"`
}

type calendarAttendee struct {
	Email          string `json:"email"`
	Self           bool   `json:"self"`
	ResponseStatus string `json:"responseStatus"`
}

func parseCalendarEvents(items []calendarEventItem) []CalendarEvent {
	var events []CalendarEvent
	for _, item := range items {
		event := CalendarEvent{
			Title:       item.Summary,
			MeetingLink: extractMeetingLink(item),
			RSVP:        extractRSVP(item.Attendees),
		}

		if item.Start.Date != "" {
			// All-day event
			event.AllDay = true
		} else if item.Start.DateTime != "" {
			t, err := time.Parse(time.RFC3339, item.Start.DateTime)
			if err == nil {
				event.StartTime = t
			}
		}

		events = append(events, event)
	}
	return events
}

func extractMeetingLink(item calendarEventItem) string {
	// Prefer conference data entry points
	if item.ConferenceData != nil {
		for _, ep := range item.ConferenceData.EntryPoints {
			if ep.EntryPointType == "video" && ep.URI != "" {
				return ep.URI
			}
		}
	}

	// Fall back to hangout link
	if item.HangoutLink != "" {
		return item.HangoutLink
	}

	return ""
}

func extractRSVP(attendees []calendarAttendee) RSVPStatus {
	for _, a := range attendees {
		if a.Self {
			switch a.ResponseStatus {
			case "accepted":
				return RSVPAccepted
			case "declined":
				return RSVPDeclined
			case "needsAction":
				return RSVPNeedsAction
			case "tentative":
				return RSVPTentative
			}
		}
	}
	// No self attendee found — treat as accepted (e.g., events the user owns)
	return RSVPAccepted
}

// ExtractMeetingLinkFromDescription looks for URLs in event descriptions
// that look like meeting links (Zoom, Teams, etc.).
func extractMeetingLinkFromDescription(description string) string {
	for _, prefix := range []string{"https://zoom.us/", "https://teams.microsoft.com/", "https://meet.google.com/"} {
		if idx := strings.Index(description, prefix); idx >= 0 {
			end := strings.IndexAny(description[idx:], " \n\r\t")
			if end == -1 {
				return description[idx:]
			}
			return description[idx : idx+end]
		}
	}
	return ""
}
