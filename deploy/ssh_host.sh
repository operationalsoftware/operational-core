#!/bin/bash

# Exit if any command fails
set -e

# Get the directory of this script
this_dir="$(cd "$(dirname "$0")" && pwd)"

# Set the working directory to that of this script
cd "$this_dir"

# Check if environment is provided
if [ -z "$1" ]; then
  echo "Error: No deployment environment specified."
  echo "Usage: $0 <environment> [ssh-args...]"
  exit 1
fi

deployment_env="$1"
shift  # remove the first argument so $@ contains only SSH args

# Source the get-deploy-config script to read the configuration values
source ./read_config.sh "$deployment_env"

# Execute the SSH command with the optional SSH key and forward extra args
ssh $DEPLOY_SSH_KEY_FLAG "$DEPLOY_HOST" "$@"

