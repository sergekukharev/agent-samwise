# Sam (agent-samwise)

A personal assistant agent inspired by Samwise Gamgee — the helper who carries you when you can't carry yourself.

## Project Overview

Sam is a Go-based personal assistant that runs as a scheduled GitHub Action and locally via CLI. It integrates with Google Calendar, Gmail, Todoist, and Slack to help manage daily priorities.

## Build & Test

```bash
go build ./...
go test ./...
```

## Code Style

- Follow the parent workspace AGENTS.md guidelines
- Domain-expressive naming — no Manager/Handler/Service/Utils
- Value objects over primitives
- No setters — use intent-capturing methods
- Interfaces for external API clients (testability)

## Configuration

- Runtime config: YAML file (calendars, projects, boards to watch)
- Secrets: environment variables (API keys, tokens)

## Git

Use Conventional Commits. Add AI co-authorship trailer.
