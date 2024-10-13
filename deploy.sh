#!/bin/bash

# Check if there are uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
  echo "Error: There are uncommitted changes. Please commit or stash your changes before deploying."
  exit 1
fi

# Check if there are any unpushed commits
if [ -n "$(git log origin/$(git rev-parse --abbrev-ref HEAD)..HEAD)" ]; then
  echo "Error: There are unpushed commits. Please push your changes before deploying."
  exit 1
fi

# Exit if any command fails
set -e

# Source the get-deploy-config script to read the configuration values
source ./get-deploy-config

# Perform the deployment using the host value
echo "Deploying branch '$(git rev-parse --abbrev-ref HEAD)' to $host..."

# Build the Golang binary
./build.sh app linux/amd64

# Copy the necessary files to the host
scp $ssh_key_flag ./app "$host:/app/app.new"
scp $ssh_key_flag ./Caddyfile "$host:/app/Caddyfile"

# Copy the services to the systemd folder (convenience)
scp $ssh_key_flag ./caddy.service "$host:/app/caddy.service"
scp $ssh_key_flag ./app.service "$host:/app/app.service"

# Rename the binaries on the host
if ssh $ssh_key_flag "$host" "[ -f /app/app ]"; then
  ssh $ssh_key_flag "$host" "mv /app/app /app/app.old"
fi
ssh $ssh_key_flag "$host" "mv /app/app.new /app/app"

# Restart the app.service on the host
ssh $ssh_key_flag "$host" "service app restart"

# Restart the caddy.service on the host
ssh $ssh_key_flag "$host" "service caddy restart"

# Remove the old binary on the host if it exists
if ssh $ssh_key_flag "$host" "[ -f /app/app.old ]"; then
  ssh $ssh_key_flag "$host" "rm /app/app.old"
fi

# Remove the local binary
rm ./app

echo "Deployment completed successfully."
