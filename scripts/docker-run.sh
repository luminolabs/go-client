#!/bin/bash

# Collect CLI arguments
args=("$@")

# Export .env environment variables
set -o allexport
eval $(cat ~/.lumino/.env | grep -v '^#' | tr -d '\r')

# If PZ_DEVICE is `cuda`, set the GPU flag
GPUS=""
if [ "$PZ_DEVICE" == "cuda" ]; then
  GPUS="--gpus all"
fi

# Run the go-client Docker container with the specified arguments
# Mount a volume to the current directory .workdir
# Set the PZ_ROOT_DIR environment variable to /pipeline-zen-jobs
docker run $GPUS -e PZ_ROOT_DIR=/pipeline-zen-jobs -e PZ_ENV_DIR=/root/.lumino -v ~/.lumino:/root/.lumino --rm -it go-client:latest "${args[@]}"