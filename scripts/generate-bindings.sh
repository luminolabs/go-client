#!/bin/bash

set -e -o pipefail

mkdir -p ./pkg/bindings

generate_binding() {
  contract=$(echo $1 | awk '{print $1}')
  go_source=$(echo $1 | awk '{print $2}')
  echo "Generating binding for (${contract})"
  abigen --abi ./abis/${contract}-abi.json --pkg 'bindings' --type=${contract} --out ./pkg/bindings/${go_source}
}

contracts=(
  # TODO: Add files as we test
  "StateManager stateManager.go"
)

for c in "${contracts[@]}"
do
    generate_binding "$c"
done

bash ./scripts/copy-contracts-address.sh