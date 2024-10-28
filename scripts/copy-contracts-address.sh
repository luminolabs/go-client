#!/bin/bash

echo "Copying addresses.json to core/contracts.go"
touch core/contracts.go

# Write the package declaration to contracts.go
echo "package core" > core/contracts.go
echo "" >> core/contracts.go  # Add an empty line for separation

# Function to generate Go variable declarations from JSON keys and addresses
generate_contract_address() {
  jsonFileKey=".$(awk '{print $1}' <<< "$1")"  # Extract the first part as JSON key
  goContractKey=$(awk '{print $2}' <<< "$1")  # Extract the second part as Go variable name
  varDeclaration="var $goContractKey ="

  # Use jq to extract the address value without quotes
  contractAddress=$(jq -r "$jsonFileKey" addresses.json)

  # Append the variable declaration with the extracted address to contracts.go
  echo "$varDeclaration \"$contractAddress\"" >> core/contracts.go
}

# List of contract keys and their corresponding Go variable names
contract_addresses_list=(
  "StateManager StateManagerAddress"
  "StakeManager StakeManagerAddress"
  "JobManager JobManagerAddress"
  "VoteManager VoteManagerAddress"
  "BlockManager BlockManagerAddress"
)

# Iterate over the list and generate the corresponding Go variables
for c in "${contract_addresses_list[@]}"; do
    generate_contract_address "$c"
done