# Go-node: Developer Documentation

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Getting Started](#2-getting-started)
3. [Project Structure](#3-project-structure)
4. [Core Components](#4-core-components)
5. [Command Line Interface](#5-command-line-interface)
6. [Development Workflow](#6-development-workflow)
7. [Testing](#7-testing)
8. [Common Patterns and Best Practices](#8-common-patterns-and-best-practices)
9. [Troubleshooting](#9-troubleshooting)

## 1. Project Overview

The Lumino Go Client is a command-line interface (CLI) application for interacting with the Lumino network. It provides functionalities for staking, job management, block operations, and network status queries.

### Key Features

- Staking and unstaking operations
- Job creation, listing, and execution
- Block proposal and confirmation
- Network and account status queries

## 2. Getting Started

### Prerequisites

- Go 1.22 or later
- Git

### Setting Up the Development Environment

1. Clone the repository:
   ```
   git clone https://github.com/your-org/lumino-go-client.git
   cd lumino-go-client
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Build the project:
   ```
   go build -o luminocli
   ```

4. Run the CLI:
   ```
   ./luminocli --help
   ```
   
## 3. Running with Docker
First, build the Docker image:
```
./scripts/docker-build.sh
```

Then, run the Lumino Client with Docker; for example, to stake 1 token:
```
./scripts/docker-run.sh ./lumino stake --address 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --value 1  --logLevel debug 
```

## 4. Project Structure

```
luminoclient/
├── cmd/                 # Command implementations
│   ├── root.go          # Root command and entry point
│   ├── stake.go         # Staking commands
│   ├── unstake.go       # Unstaking commands
│   ├── job.go           # Job-related commands
│   ├── block.go         # Block operations commands
│   ├── status.go        # Status query commands
│   ├── config.go        # Configuration management
│   ├── utils.go         # Command utilities
│   └── interface.go     # Interfaces for components
├── core/                # Core types and constants
│   ├── constants.go     # System-wide constants
│   ├── contracts.go     # Contract addresses
│   ├── version.go       # Version information
│   └── types/           # Data structures
│       ├── block.go
│       ├── epoch.go
│       ├── job.go
│       └── staker.go
├── logger/              # Logging functionality
│   ├── logger.go
│   └── errors.go
├── utils/               # Utility functions
│   ├── api.go
│   ├── signature.go
│   ├── transaction.go
│   └── utils.go
├── go.mod
├── go.sum
└── main.go              # Application entry point
```

## 5. Core Components

### 5.1 Command (cmd) Package

The `cmd` package is the heart of the CLI, implementing various commands using the Cobra library.

Key concepts:
- Each command is defined as a Cobra command struct
- Commands are organized hierarchically (root command -> subcommands)
- Command execution logic is defined in the `Run` field of each command

Example of adding a new command:

```go
var newCmd = &cobra.Command{
    Use:   "new",
    Short: "A new command",
    Run: func(cmd *cobra.Command, args []string) {
        // Command logic here
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

### 5.2 Core Types

The `core/types` package defines the main data structures used throughout the application. When working with these types, ensure you understand their relationships and how they map to the Lumino network concepts.

### 5.3 Utility Functions

The `utils` package contains helper functions for common operations. When adding new functionality, consider if it can be generalized and added to this package for reuse.

## 6. Command Line Interface

The CLI is built using the Cobra library. Key points to remember:

- The root command is defined in `cmd/root.go`
- Each major feature has its own command file (e.g., `stake.go`, `job.go`)
- Use flags for command options (defined using `cmd.Flags().StringP()` or similar methods)
- Implement `Run` functions for each command to define its behavior

## 7. Development Workflow

1. **Feature Planning**: Discuss new features in the issue tracker before implementation.
2. **Branch Creation**: Create a new branch for each feature or bug fix.
3. **Implementation**: Write code and tests for the new feature.
4. **Testing**: Run tests and ensure all existing tests pass.
5. **Documentation**: Update relevant documentation, including this developer guide if necessary.
6. **Pull Request**: Create a pull request for code review.
7. **Code Review**: Address any feedback from the code review.
8. **Merge**: Once approved, merge the pull request into the main branch.

## 8. Testing

- Write unit tests for all new functionality
- Use table-driven tests for testing multiple scenarios
- Mock external dependencies (e.g., blockchain interactions) for isolated testing
- Run tests using `go test ./...` from the project root

Example of a table-driven test:

```go
func TestSomeFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case 1", "input1", "expected1"},
        {"case 2", "input2", "expected2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := SomeFunction(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## 9. Common Patterns and Best Practices

- Use interfaces for better testability and modularity (see `cmd/interface.go`)
- Follow Go naming conventions (e.g., use MixedCaps or mixedCaps)
- Handle errors explicitly and avoid using panic
- Use context for managing timeouts and cancellations in long-running operations
- Prefer composition over inheritance

## 10. Troubleshooting

Common issues and their solutions:

1. **Build Errors**: 
   - Ensure all dependencies are installed (`go mod tidy`)
   - Check for conflicting versions in `go.mod`

2. **Runtime Errors**:
   - Check log files for detailed error messages
   - Ensure configuration files are correctly set up

3. **Test Failures**:
   - Run tests with verbose output (`go test -v ./...`)
   - Check for race conditions with (`go test -race ./...`)

For any other issues, consult the project's issue tracker or reach out to the core development team.

This documentation is a living document. As we work on the project and gain new insights, feel free to update and expand this guide to help future developers.