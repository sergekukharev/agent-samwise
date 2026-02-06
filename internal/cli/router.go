package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/sergekukharev/agent-samwise/internal/config"
	"github.com/sergekukharev/agent-samwise/internal/output"
)

// Router dispatches subcommands to registered capabilities.
type Router struct {
	capabilities map[string]Capability
}

// NewRouter creates a Router with the given capabilities.
func NewRouter(capabilities []Capability) *Router {
	m := make(map[string]Capability, len(capabilities))
	for _, c := range capabilities {
		m[c.Name] = c
	}
	return &Router{capabilities: m}
}

// Run parses arguments and dispatches to the appropriate capability.
// It returns an exit code (0 for success, 1 for errors).
func (r *Router) Run(args []string) int {
	fs := flag.NewFlagSet("sam", flag.ContinueOnError)
	configPath := fs.String("config", config.DefaultPath, "path to config file")

	// Parse global flags from the full argument list.
	// flag.FlagSet stops at the first non-flag argument (the subcommand).
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	remaining := fs.Args()

	if len(remaining) == 0 {
		r.printHelp()
		return 0
	}

	subcmd := remaining[0]
	cap, ok := r.capabilities[subcmd]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", subcmd)
		r.printSuggestions(subcmd)
		return 1
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	for _, name := range cap.RequiredConfig {
		if err := cfg.ValidateFor(name); err != nil {
			fmt.Fprintf(os.Stderr, "config error: %v\n", err)
			return 1
		}
	}

	secrets, err := config.ResolveSecrets(cap.RequiredEnv...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	presenter, err := output.DetectPresenter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if err := cap.Run(cfg, secrets, presenter); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func (r *Router) printHelp() {
	fmt.Println("Sam — your personal assistant")
	fmt.Println()
	fmt.Println("Usage: sam [flags] <command>")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --config <path>  path to config file (default: config.yaml)")
	fmt.Println()
	fmt.Println("Commands:")

	names := r.sortedNames()
	if len(names) == 0 {
		fmt.Println("  (no capabilities registered)")
		return
	}

	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	for _, name := range names {
		cap := r.capabilities[name]
		fmt.Printf("  %-*s  %s\n", maxLen, name, cap.Description)
	}
}

func (r *Router) printSuggestions(unknown string) {
	names := r.sortedNames()
	var suggestions []string
	for _, name := range names {
		if strings.HasPrefix(name, unknown[:1]) || levenshtein(unknown, name) <= 3 {
			suggestions = append(suggestions, name)
		}
	}

	if len(suggestions) > 0 {
		fmt.Fprintf(os.Stderr, "Did you mean:\n")
		for _, s := range suggestions {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
		fmt.Fprintln(os.Stderr)
	}

	fmt.Fprintf(os.Stderr, "Run 'sam' for a list of available commands.\n")
}

func (r *Router) sortedNames() []string {
	names := make([]string, 0, len(r.capabilities))
	for name := range r.capabilities {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
		d[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			d[i][j] = min(d[i-1][j]+1, min(d[i][j-1]+1, d[i-1][j-1]+cost))
		}
	}
	return d[la][lb]
}
