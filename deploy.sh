#!/bin/bash

# Check if there are uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
  echo "Error: There are uncommitted changes. Please commit or stash your changes before deploying."
  exit 1
fi

# Exit if any command fails
set -e

# Read the configuration file based on the environment (current git branch)
config_file="deploy.config"
host=$(awk -F "=" -v env="$(git rev-parse --abbrev-ref HEAD)" '$1 == "[" env "]" {f=1; next} f && /^Host/ {print $2; exit}' "$config_file")
ssh_key=$(awk -F "=" -v env="$(git rev-parse --abbrev-ref HEAD)" '$1 == "[" env "]" {f=1; next} f && /^SSHKey/ {print $2; exit}' "$config_file")

# Validate the host value
if [ -z "$host" ]; then
  echo "Error: Host not found for current git branch in $config_file"
  exit 1
fi

# Perform the deployment using the host value
echo "Deploying branch '$(git rev-parse --abbrev-ref HEAD)' to $host..."

# Build the Golang binary
./build.sh app linux/amd64

# Define the scp and ssh commands with or without SSHKey parameter
if [ -n "$ssh_key" ]; then
  scp_cmd="scp -i $ssh_key"
  ssh_cmd="ssh -i $ssh_key"
else
  scp_cmd="scp"
  ssh_cmd="ssh"
fi

# Copy the necessary files to the host
$scp_cmd ./app "$host:/app/app.new"
$scp_cmd ./Caddyfile "$host:/app/Caddyfile"

# the following copies are for convenience: the services need to be copied to
# the systemd folder and reloaded manually
$scp_cmd ./caddy.service "$host:/app/caddy.service"
$scp_cmd ./app.service "$host:/app/app.service"

# Rename the binaries on the host
# check if /app/app exists and rename it to /app/app.old
if $ssh_cmd $host "[ -f /app/app ]"; then
  $ssh_cmd $host "mv /app/app /app/app.old"
fi
$ssh_cmd $host "mv /app/app.new /app/app"

# Restart the app.service on the host
$ssh_cmd $host "service app restart"

# Restart the caddy.service on the host
$ssh_cmd $host "service caddy restart"

# Remove the old binary on the host if it exists
if $ssh_cmd $host "[ -f /app/app.old ]"; then
  $ssh_cmd $host "rm /app/app.old"
fi

# Remove the local binary
rm ./app

echo "Deployment completed successfully."
