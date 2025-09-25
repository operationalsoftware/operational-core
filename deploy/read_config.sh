#!/bin/bash

config_file="deploy.config"

# Check if environment name is provided
if [ -z "$1" ]; then
  echo "Error: No deployment environment specified."
  echo "Usage: $0 <environment>"
  exit 1
fi

deployment_env="$1"

# Check if environment exists in config
if ! grep -q "^\[$deployment_env\]" "$config_file"; then
  echo "Error: Environment '$deployment_env' not found in $config_file"
  exit 1
fi

# Helper to read and trim a key from the config
read_config_value() {
  local key="$1"
  awk -F "=" -v deployment_env="[$deployment_env]" -v key="$key" '
    $1 == deployment_env {f=1; next}
    f && $1 ~ "^"key {
      val=$2
      gsub(/^[ \t]+|[ \t]+$/, "", val)  # trim leading/trailing spaces
      print val
      exit
    }
  ' "$config_file"
}

# Extract values
host=$(read_config_value "Host")
ssh_key=$(read_config_value "SSHKey")
requires_confirmation=$(read_config_value "RequiresConfirmation")

# Validate the extracted values
if [ -z "$host" ]; then
  echo "Error: Host not found for deployment environment '$deployment_env' in $config_file"
  exit 1
fi

# Default requires_confirmation to false if not set
if [ -z "$requires_confirmation" ]; then
  requires_confirmation="false"
fi

# Prepare the SSH key flag
if [ -n "$ssh_key" ]; then
  ssh_key_flag="-i $ssh_key"
else
  ssh_key_flag=""
fi

# Export the results
export DEPLOY_HOST="$host"
export DEPLOY_SSH_KEY_FLAG="$ssh_key_flag"
export DEPLOY_REQUIRES_CONFIRMATION="$requires_confirmation"

# If running standalone, also print them
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  echo "Host: $DEPLOY_HOST"
  echo "SSH Key Flag: $DEPLOY_SSH_KEY_FLAG"
  echo "Requires Confirmation: $DEPLOY_REQUIRES_CONFIRMATION"
fi
