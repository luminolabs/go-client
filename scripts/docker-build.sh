#!/bin/bash

# Check if `pipeline-zen-src` directory exists
# If it does, git pull to update the repository
# else, clone the repository
if [ -d "pipeline-zen-src" ]; then
  echo "Updating pipeline-zen repository"
  cd pipeline-zen-src
  git pull
  cd ..
else
  echo "Cloning pipeline-zen repository"
  git clone git@github.com:luminolabs/pipeline-zen.git pipeline-zen-src
fi

# Build the Docker image
echo "Building the Docker image; tag: go-client:latest"
docker buildx build . -t go-client:latest