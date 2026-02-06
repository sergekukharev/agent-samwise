# cli Specification

## Purpose
TBD - created by archiving change project-foundation. Update Purpose after archive.
## Requirements
### Requirement: Subcommand-based CLI
The system SHALL expose each capability as a separate subcommand (e.g., `sam calendar-sync`, `sam review-projects`).

#### Scenario: Running a known subcommand
- **WHEN** user runs `sam <subcommand>`
- **THEN** the corresponding capability executes with loaded configuration

#### Scenario: Running without a subcommand
- **WHEN** user runs `sam` with no arguments
- **THEN** help text listing all available subcommands is displayed

#### Scenario: Running an unknown subcommand
- **WHEN** user runs `sam <unknown>`
- **THEN** an error message is displayed with suggestions for valid subcommands

### Requirement: Global flags
The system SHALL support global flags that apply to all subcommands.

#### Scenario: Config file override
- **WHEN** user runs `sam --config /path/to/config.yaml <subcommand>`
- **THEN** the specified config file is used instead of the default location

