# output-delivery Specification

## Purpose
TBD - created by archiving change project-foundation. Update Purpose after archive.
## Requirements
### Requirement: Context-aware output routing
The system SHALL route output to Slack or terminal based on the execution context.

#### Scenario: Running in GitHub Actions
- **WHEN** the application detects it is running in GitHub Actions (via `GITHUB_ACTIONS` environment variable)
- **AND** `SLACK_WEBHOOK_URL` environment variable is set
- **THEN** output is delivered as a Slack message

#### Scenario: Running locally
- **WHEN** the application is running outside GitHub Actions
- **THEN** output is printed to stdout in a human-readable format

#### Scenario: GitHub Actions without Slack webhook
- **WHEN** the application is running in GitHub Actions
- **AND** `SLACK_WEBHOOK_URL` is not set
- **THEN** the application exits with a non-zero exit code
- **AND** an error message indicates the missing webhook URL

### Requirement: Structured output content
The system SHALL format output as markdown with a title and sections.

#### Scenario: Slack message formatting
- **WHEN** output is delivered to Slack
- **THEN** the message uses Slack mrkdwn formatting in a single section block

#### Scenario: Terminal output formatting
- **WHEN** output is printed to terminal
- **THEN** the output uses markdown with `#` for title and `##` for section headings

