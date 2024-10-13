#!/bin/bash

# Read the configuration file based on the environment (current git branch)
config_file="deploy.config"
branch="$(git rev-parse --abbrev-ref HEAD)"

# Extract host and SSH key from the config file
host=$(awk -F "=" -v env="[$branch]" '$1 == env {f=1; next} f && /^Host/ {print $2; exit}' "$config_file")
ssh_key=$(awk -F "=" -v env="[$branch]" '$1 == env {f=1; next} f && /^SSHKey/ {print $2; exit}' "$config_file")

# Validate the extracted values
if [ -z "$host" ]; then
  echo "Error: Host not found for branch '$branch' in $config_file"
  exit 1
fi

# Prepare the SSH key flag
if [ -n "$ssh_key" ]; then
  ssh_key_flag="-i $ssh_key"
else
  ssh_key_flag=""
fi

# Output the results
echo "Host: $host"
echo "SSH Key Flag: $ssh_key_flag"

