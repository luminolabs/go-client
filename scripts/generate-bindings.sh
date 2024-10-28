#!/bin/bash

# Set the script to exit immediately if any commands return a non-zero status
set -e -o pipefail

# Create the bindings directory if it doesn't exist
mkdir -p ./pkg/bindings

# Function to generate Go bindings for a contract
generate_binding() {
  contract=$(echo $1 | awk '{print $1}')
  go_source=$(echo $1 | awk '{print $2}')
  echo "Generating binding for (${contract})"
  abigen --abi ./abis/${contract}-abi.json --pkg 'bindings' --type=${contract} --out ./pkg/bindings/${go_source}
}

# List of contracts and their corresponding Go source files
contracts=(
  # TODO: Add files as we test
  "StateManager stateManager.go"
  "StakeManager stakeManager.go"
  "BlockManager blockManager.go"
  "JobManager jobManager.go"
)

# Generate bindings for each contract
for c in "${contracts[@]}"
do
    generate_binding "$c"
done

# Copy the contract addresses to the Go file
./scripts/copy-contracts-address.sh