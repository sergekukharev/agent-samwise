## ADDED Requirements

### Requirement: YAML configuration loading
The system SHALL load configuration from a YAML file specifying which resources to watch and act on.

#### Scenario: Loading a valid config file
- **WHEN** the application starts with a valid YAML config file
- **THEN** all configured calendars, projects, boards, and areas are available to capabilities

#### Scenario: Config file not found
- **WHEN** the application starts and the config file does not exist at the expected path
- **THEN** the application exits with a clear error message indicating the expected config location

#### Scenario: Invalid YAML syntax
- **WHEN** the application starts with a malformed YAML config file
- **THEN** the application exits with an error message including the parse error and line number

### Requirement: Default config location
The system SHALL look for configuration at `config.yaml` in the current working directory by default.

#### Scenario: Default path resolution
- **WHEN** no `--config` flag is provided
- **THEN** the application loads configuration from `./config.yaml`

### Requirement: Environment variable resolution for secrets
The system SHALL resolve API keys and tokens from environment variables, not from the config file.

#### Scenario: All required env vars present
- **WHEN** all required environment variables for a capability are set
- **THEN** the capability has access to valid credentials

#### Scenario: Missing required env var
- **WHEN** a required environment variable is not set
- **THEN** the application exits with an error naming the missing variable and which capability needs it

### Requirement: Config validation
The system SHALL validate configuration completeness for each capability before execution.

#### Scenario: Capability-specific validation
- **WHEN** a subcommand is invoked
- **THEN** only the configuration sections required by that capability are validated
- **AND** missing optional sections for other capabilities do not cause errors
