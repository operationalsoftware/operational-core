#!/bin/bash

# Usage: ./ssh-host.sh [branch]

config_file="deploy.config"
branch=$1

# If branch is not provided, retrieve the current git branch
if [ -z "$branch" ]; then
  branch=$(git rev-parse --abbrev-ref HEAD)
fi

host_config=$(awk -F "=" -v env="$branch" '$1 == "[" env "]" {f=1; next} f && /^Host/ {print $2; exit}' "$config_file")
ssh_key=$(awk -F "=" -v env="$branch" '$1 == "[" env "]" {f=1; next} f && /^SSHKey/ {print $2; exit}' "$config_file")

# Validate the host_config value
if [ -z "$host_config" ]; then
  echo "Error: Host configuration not found for branch '$branch' in $config_file"
  exit 1
fi

# Extract the host and other parameters from the host_config
IFS=" " read -ra host_params <<< "$host_config"
host="${host_params[0]}"
ssh_key="${host_params[1]}"

# Validate the extracted values
if [ -z "$host" ]; then
  echo "Error: Host not found for branch '$branch' in $config_file"
  exit 1
fi

echo "SSH into host for branch '$branch'..."

# SSH into the host
if [ -n "$ssh_key" ]; then
  ssh -i "$ssh_key" "$host"
else
  ssh "$host"
fi

