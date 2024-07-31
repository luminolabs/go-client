# Lumino Node Client

This is the Lumino Node Client, a Go application for participating in the Lumino network as a compute provider.

## Getting Started

### Prerequisites

- Go 1.22.5 or later

### Installation

1. Clone the repository
2. Run `go mod tidy` to ensure all dependencies are downloaded
3. Build the application with `go build -o lumino-client`

### Running the Client

Run the compiled binary:

```
./lumino-client
```

## Project Structure

- `main.go`: Entry point of the application
- `client/`: Main client logic
- `config/`: Configuration management
- `ml/`: Machine Learning container
- `blockchain/`: Blockchain interaction container
- `monitoring/`: Monitoring service
- `job/`: Job management
- `dashboard/`: Web dashboard for the client

## Contributing

[Add contribution guidelines here]

## License

[Add license information here]