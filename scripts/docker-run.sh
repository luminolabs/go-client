#!/bin/bash

# Collect CLI arguments
args=("$@")

# Run the go-client Docker container with the specified arguments
docker run --rm -it go-client:latest "${args[@]}"