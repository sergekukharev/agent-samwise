## Context

This is a greenfield Go project. We need a CLI skeleton, config system, and output routing before building any capabilities. The foundation must support two execution contexts: local interactive use and scheduled GitHub Actions runs.

## Goals / Non-Goals

**Goals:**
- Establish project layout and Go module
- CLI with subcommand routing and global flags
- YAML config loading with per-capability validation
- Output system that auto-detects Slack vs terminal
- GitHub Actions workflow for daily scheduled runs
- Interfaces for external APIs to enable testing

**Non-Goals:**
- Implementing any actual capabilities (calendar sync, email triage, etc.)
- OAuth flows — credentials come from environment variables
- Persistent storage — Sam is stateless between runs
- User-facing API or web interface

## Decisions

### Decision: Standard library for CLI
Using Go's standard library (`os.Args` + `flag`) for CLI routing and flag parsing.

**Rationale:** Sam's CLI needs are simple — a handful of subcommands with shared global flags. A lightweight subcommand router (map of name → function) keeps dependencies minimal and the code transparent. No framework magic.

**Alternatives considered:**
- `github.com/spf13/cobra` — Powerful but heavyweight for a tool with ~5 subcommands. Adds a dependency tree we don't need.
- `urfave/cli` — Same concern — unnecessary dependency for simple routing

### Decision: gopkg.in/yaml.v3 for config parsing
Using `gopkg.in/yaml.v3` for YAML configuration.

**Rationale:** Standard Go YAML library with good error messages including line numbers. Supports struct tags for mapping.

**Alternatives considered:**
- `github.com/goccy/go-yaml` — Faster but less mature, error messages not as good
- TOML/JSON — YAML is more human-friendly for a config file users will edit

### Decision: Slack Incoming Webhooks for output
Using Slack Incoming Webhooks (not the Slack API) for message delivery.

**Rationale:** Simplest integration — a single URL, no OAuth, no bot token management. A webhook URL as an env var is the lowest-friction approach for a personal tool.

**Alternatives considered:**
- Slack Bot API — More powerful (threads, reactions) but requires OAuth setup, bot installation. Overkill for one-way notifications.
- Email — Less immediate than Slack, user already wants Slack

### Decision: Project layout
```
.
├── cmd/
│   └── sam/
│       └── main.go          # Entry point
├── internal/
│   ├── cli/                  # Cobra command definitions
│   ├── config/               # Config loading and validation
│   ├── output/               # Output routing (Slack/terminal)
│   └── platform/             # External API client interfaces
├── .github/
│   └── workflows/
│       └── daily.yml         # Scheduled GitHub Action
├── go.mod
└── go.sum
```

**Rationale:** Standard Go project layout. `internal/` prevents external imports. `cmd/sam/` allows the binary name to be `sam`. Each concern gets its own package. `platform/` holds interfaces — concrete implementations will be added when capabilities are built.

### Decision: Output as an interface
The output system will be an interface, not concrete types:

```go
type Briefing struct {
    Title    string
    Sections []Section
}

type Section struct {
    Heading string
    Body    string
}

type Presenter interface {
    Present(briefing Briefing) error
}
```

**Rationale:** Capabilities produce `Briefing` values. The presenter (Slack or terminal) is injected based on context. This keeps capabilities decoupled from delivery and makes testing trivial — just assert on the Briefing content.

### Decision: Capability registration pattern
Each capability registers as a named command with a `Run` function that accepts config and a `Presenter`:

```go
type Capability struct {
    Name        string
    Description string
    Run         func(cfg config.Config, out output.Presenter) error
    RequiredEnv []string
}
```

The CLI main function maintains a map of `name → Capability`. Dispatching is `capabilities[os.Args[1]].Run(cfg, presenter)`.

**Rationale:** Keeps capabilities self-contained. Adding a new capability means creating a new `Capability` value and adding it to the map — no framework, no code generation.

## Risks / Trade-offs

- **Slack webhooks are one-way** → Cannot thread messages or receive user input. Acceptable for a daily briefing tool. If interactivity is needed later, migrate to Bot API.
- **Missing Slack webhook in CI = hard failure** → In GitHub Actions, if `SLACK_WEBHOOK_URL` is not set, the build fails. This is intentional — the user gets a GitHub notification email and knows something is wrong, rather than silently losing output.
- **No authentication flow** → User must manually obtain API tokens and set env vars. Acceptable for a personal tool. Could add a `sam auth` setup wizard later.
- **Standard library CLI is more manual** → Requires writing help text and flag parsing by hand. Acceptable for ~5 subcommands — simplicity outweighs convenience.

## Open Questions

- Should the daily GitHub Action run all capabilities sequentially, or should each capability be a separate workflow/job? Starting with sequential (single workflow, multiple subcommand invocations) for simplicity.
