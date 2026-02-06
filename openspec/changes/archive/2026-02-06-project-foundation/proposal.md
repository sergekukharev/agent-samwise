## Why

Sam has five planned capabilities (calendar sync, project review, email triage, memory, calendar recommendations), and all of them need shared infrastructure: a CLI entry point, configuration loading, output delivery (Slack vs terminal), and a GitHub Actions workflow. Without this foundation, each capability would reinvent these concerns independently.

## What Changes

- Initialize Go module with CLI skeleton supporting subcommands per capability
- YAML-based configuration loader for calendars, projects, boards, and other user-specific settings
- Output delivery system that routes messages to Slack (scheduled runs) or terminal (local runs), determined by execution context
- GitHub Actions workflow for scheduled daily execution
- Shared types and interfaces for external API clients (Google, Todoist, Slack)

## Capabilities

### New Capabilities
- `cli`: Command-line interface structure — subcommand routing, flag parsing, help text
- `configuration`: YAML config loading, validation, and environment variable resolution for secrets
- `output-delivery`: Routing output to Slack or terminal based on execution context

### Modified Capabilities
<!-- None — this is a greenfield project -->

## Impact

- New Go module (`github.com/sergekukharev/agent-samwise`)
- New dependencies: YAML parsing, Slack client, CLI framework
- New GitHub Actions workflow file (`.github/workflows/daily.yml`)
- Establishes project layout and patterns all future capabilities will follow
