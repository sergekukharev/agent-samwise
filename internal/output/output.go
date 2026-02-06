package output

// Briefing is the structured output that capabilities produce.
type Briefing struct {
	Title    string
	Sections []Section
}

// Section is a titled block of content within a Briefing.
type Section struct {
	Heading string
	Body    string
}

// Presenter delivers a Briefing to the user.
type Presenter interface {
	Present(briefing Briefing) error
}
