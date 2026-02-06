# Tasks

## 1. Go Module & Project Layout
- [x] 1.1 Initialize Go module (`go mod init github.com/sergekukharev/agent-samwise`)
- [x] 1.2 Create directory structure: `cmd/sam/`, `internal/cli/`, `internal/config/`, `internal/output/`, `internal/platform/`
- [x] 1.3 Create `cmd/sam/main.go` with minimal entry point that prints help

## 2. Configuration
- [x] 2.1 Define `Config` struct in `internal/config/` with sections for each planned capability (calendars, todoist, gmail, slack, areas of interest)
- [x] 2.2 Implement YAML config loader using `gopkg.in/yaml.v3` with error messages including line numbers
- [x] 2.3 Implement default path resolution (`./config.yaml` in working directory)
- [x] 2.4 Implement `--config` flag override
- [x] 2.5 Implement per-capability config validation (only validate sections needed by the invoked subcommand)
- [x] 2.6 Implement environment variable resolution for secrets with clear error on missing vars
- [x] 2.7 Write tests for config loading, validation, and env var resolution

## 3. Output Delivery
- [x] 3.1 Define `Briefing`, `Section`, and `Presenter` interface in `internal/output/`
- [x] 3.2 Implement `TerminalPresenter` — plain text with headers and indentation
- [x] 3.3 Implement `SlackPresenter` — Slack Block Kit formatting via incoming webhook
- [x] 3.4 Implement presenter auto-detection: in GitHub Actions, require `SLACK_WEBHOOK_URL` (exit non-zero if missing); locally, use terminal
- [x] 3.5 Write tests for terminal presenter output and Slack payload construction

## 4. CLI Routing
- [x] 4.1 Implement subcommand routing using `os.Args` and a capability registry map
- [x] 4.2 Implement global flags (`--config`) using `flag` package
- [x] 4.3 Implement help text generation from registered capabilities
- [x] 4.4 Implement unknown subcommand error with suggestions
- [x] 4.5 Wire config loading → presenter detection → capability dispatch in main
- [x] 4.7 Write tests for subcommand routing and flag parsing

## 5. GitHub Actions
- [x] 5.1 Create `.github/workflows/daily.yml` with cron schedule for morning run
- [x] 5.2 Configure workflow to build and run `sam` with env vars from GitHub Secrets
- [x] 5.3 Add manual trigger (`workflow_dispatch`) for on-demand runs

## 6. Sample Capability (placeholder)
- [x] 6.1 Create a `hello` subcommand that outputs a greeting via the presenter — proves the full pipeline works (config → capability → output)
- [ ] 6.2 Remove `hello` capability once a real capability is implemented
