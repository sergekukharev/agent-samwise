package output

import (
	"fmt"
	"io"
	"os"
)

// TerminalPresenter prints briefings to stdout as markdown.
type TerminalPresenter struct {
	Writer io.Writer
}

// NewTerminalPresenter creates a presenter that writes to stdout.
func NewTerminalPresenter() *TerminalPresenter {
	return &TerminalPresenter{Writer: os.Stdout}
}

func (p *TerminalPresenter) Present(briefing Briefing) error {
	fmt.Fprintf(p.Writer, "# %s\n\n", briefing.Title)

	for i, section := range briefing.Sections {
		fmt.Fprintf(p.Writer, "## %s\n\n", section.Heading)
		fmt.Fprintf(p.Writer, "%s\n", section.Body)
		if i < len(briefing.Sections)-1 {
			fmt.Fprintln(p.Writer)
		}
	}

	return nil
}
