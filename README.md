# Go-client: Developer Documentation

## Overview

The Go Client is a command-line interface (CLI) tool for interacting with the Lumino network, a decentralized compute platform that enables seamless execution of machine learning workflows. This client facilitates staking, job management, and ML pipeline execution in a decentralized environment.

## Table of Contents

1.  [Features](#features)
2.  [Requirements](#requirements)
3.  [Installation](#installation)
4.  [Configuration](#configuration)
5.  [Usage](#usage)
6.  [Machine Learning Pipeline Integration](#machine-learning-pipeline-integration)
7.  [Development](#development)
8.  [Testing](#testing)
9.  [Contributing](#contributing)
10. [Troubleshooting](#troubleshooting)

## Features

- **Account Management**

  - Create and import Ethereum accounts
  - Secure keystore management
  - Private key and password handling

- **Staking Operations**

  - Stake LUMINO tokens
  - Unstake tokens with lock period
  - Withdraw unlocked stakes
  - View staking status and rewards

- **Job Management**

  - Create ML training jobs
  - Assign jobs to compute providers
  - Track job status and execution
  - View job results and metrics

- **ML Pipeline Integration**

  - Seamless integration with pipeline-zen
  - Support for PyTorch training workflows
  - Multi-GPU training capabilities
  - Dataset management and versioning

- **Network Operations**
  - Monitor network status
  - Track epochs and states
  - View network metrics

## Requirements

- Go 1.22 or later
- Python 3.10 (for ML pipeline integration)
- Docker (recommended for deployment)
- Git
- Ethereum client (geth)

## Installation

### Local Installation

1. Clone the repository:

```bash
git clone https://github.com/luminolabs/go-client.git
cd go-client
```

2. Install dependencies:

```bash
go mod download
```

3. Generate contract bindings:

```bash
./scripts/generate-bindings.sh
```

4. Build the client:

```bash
go build -o lumino
```

### Docker Installation

1. Build the Docker image:

```bash
./scripts/docker-build.sh
```

2. Run the client using Docker:

```bash
./scripts/docker-run.sh ./lumino [command]
```

## Configuration

### Directory Structure

Create a `.lumino` directory in your home folder with the following structure:

```
~/.lumino/
├── .env                               # Environment variables
├── config.json                        # Configuration file
└── pipeline-zen-jobs-gcp-key.json    # GCP credentials (if using GCP)
```

### Environment Variables

Create a `.env` file with the following settings:

```dotenv
PZ_ENV=cpnode                          # Environment (cpnode/local)
PZ_RESULTS_BUCKET_SUFFIX=us            # Results bucket location
PZ_HUGGINGFACE_TOKEN=your_token_here   # HuggingFace API token
PZ_DEVICE=cuda                         # Device type (cuda/cpu)
```

### Configuration File (Optional)

Create a `config.json` file with job settings:

```json
{
  "job_config_name": "llm_dummy",
  "job_id": "13",
  "dataset_id": "gs://lum-pipeline-zen-jobs-us/datasets/your-dataset-id",
  "batch_size": "20",
  "shuffle": "true",
  "num_epochs": "1",
  "use_lora": "true",
  "use_qlora": "false",
  "lr": "1e-2",
  "override_env": "prod",
  "seed": "42",
  "num_gpus": "1",
  "user_id": "0x4118CFD00dD5e8CED96e0ff8061F56F2d155e83B"
}
```

Then, run the Lumino Client with Docker; for example, to stake 1 token:

```bash
./scripts/docker-run.sh ./lumino stake --address 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --value 1  --logLevel debug
```

Finally, run the pipeline-zen workflow with Docker:

```bash
./scripts/docker-run.sh ./lumino executeJob -a 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --config /root/.lumino/config.json --jobId 21 --zen-path /pipeline-zen-jobs --logLevel debug
```

## Usage

### Account Management

Create a new account:

```bash
./lumino create
```

Import an existing account:

```bash
./lumino import
```

### Staking Operations

Stake tokens:

```bash
./lumino stake --address <your-address> --value <amount> --logLevel debug
```

Unstake tokens:

```bash
./lumino unstake --address <your-address> --value <amount>
```

Withdraw unlocked stakes:

```bash
./lumino withdraw --address <your-address>
```

### Job Management

Create a new job:

```bash
./lumino createJob -a <your-address> --config /path/to/config.json --jobFee <amount>
```

Execute a job:

```bash
./lumino executeJob -a <your-address> --jobId <id> --zen-path /pipeline-zen-jobs --logLevel debug
```

### Network Information

View network status:

```bash
./lumino networkInfo
```

## Machine Learning Pipeline Integration

### Supported ML Tasks

- Model fine-tuning
- Distributed training
- Multi-GPU training
- Custom dataset processing

### Pipeline Configuration

The pipeline configuration supports various ML parameters:

- Batch size and learning rate
- Training epochs
- Model architecture settings
- LoRA and QLoRA support
- Dataset configurations
- GPU utilization settings

### Running ML Jobs

1. Prepare your job configuration
2. Create the job using `createJob`
3. Execute the job using `executeJob`
4. Monitor progress through logs

## Development

### Project Structure

```
lumino/
├── cmd/                # Command implementations
├── core/               # Core types and constants
├── logger/             # Logging functionality
├── utils/              # Utility functions
├── accounts/           # Account management
├── path/               # Path handling
└── pipeline-zen/       # ML pipeline integration
```

### Building from Source

1. Clone the repository
2. Install dependencies:

```bash
go mod download
```

3. Generate contract bindings:

```bash
./scripts/generate-bindings.sh
```

4. Build:

```bash
go build -o lumino
```

## Testing

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

Run all tests:

```bash
go test ./... -v
```

Run tests with race condition detection:

```bash
go test -race ./...
```

Generate test coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Contributing

1. **Feature Planning**: Discuss new features in the issue tracker before implementation.
2. **Branch Creation**: Create a new branch for each feature or bug fix.
3. **Implementation**: Write code and tests for the new feature.
4. **Testing**: Run tests and ensure all existing tests pass.
5. **Documentation**: Update relevant documentation, including this developer guide if necessary.
6. **Pull Request**: Create a pull request for code review.
7. **Code Review**: Address any feedback from the code review.
8. **Merge**: Once approved, merge the pull request into the main branch.

### Common Patterns and Best Practices

- Use interfaces for better testability and modularity (see `cmd/interface.go`)
- Follow Go naming conventions (e.g., use MixedCaps or mixedCaps)
- Handle errors explicitly and avoid using panic
- Use context for managing timeouts and cancellations in long-running operations
- Prefer composition over inheritance

## Troubleshooting

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
