#!/bin/bash

# Collect CLI arguments
args=("$@")

# Run the go-client Docker container with the specified arguments
# Mount a volume to the current directory .workdir
# Set the PZ_ROOT_DIR environment variable to /pipeline-zen-jobs
docker run -e PZ_ROOT_DIR=/pipeline-zen-jobs -e PZ_ENV_DIR=/root/.lumino -v ~/.lumino:/root/.lumino --rm -it go-client:latest "${args[@]}"