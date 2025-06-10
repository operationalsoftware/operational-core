#!/bin/bash

# Exit if any command fails
set -e

# Get the directory of this script
deploy_dir="$(cd "$(dirname "$0")" && pwd)"

# Set the working directory to that of this script
cd "$deploy_dir"

# Source the get-deploy-config script to read the configuration values
. ./get-deploy-config.sh

# Go to the main project directory
cd "$deploy_dir/.."

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

# Perform the deployment using the host value
echo "Deploying branch '$(git rev-parse --abbrev-ref HEAD)' to $host..."

# Build the Golang binary
./build.sh app linux/amd64

# Copy built app to server
scp $ssh_key_flag ./app "$host:~"
scp $ssh_key_flag "./deploy/$branch.env" "$host:~/.env"

# Remove the local binary
rm ./app

# Move back into this directory
cd "$deploy_dir"

# Copy config files to server
scp $ssh_key_flag ./Caddyfile "$host:~"
scp $ssh_key_flag ./caddy.service "$host:~"
scp $ssh_key_flag ./app.service "$host:~"
scp $ssh_key_flag ./db-backup.service "$host:~"
scp $ssh_key_flag ./db-backup.sh "$host:~"
scp $ssh_key_flag ./db-backup.timer "$host:~"

# Running commands in remote server via ssh
ssh $ssh_key_flag "$host" <<EOF
    set -e
    sudo mkdir -p /opt/app

    echo "📦 Moving config and app files..."
    sudo mv ./Caddyfile /etc/caddy/Caddyfile"
    sudo mv ./caddy.service /etc/systemd/system/caddy.service"
    sudo mv ./app.service /etc/systemd/system/app.service"
    sudo mv ./db-backup.service /etc/systemd/system/db-backup.service"
    sudo mv ./db-backup.timer /etc/systemd/system/db-backup.timer"
    sudo mv ./app /opt/app/app.new
    sudo mv "./.env" /opt/app/.env
    sudo mv ./db-backup.sh /opt/app/db-backup.sh
    sudo chmod +x /opt/app/db-backup.sh

    echo "📦 Setting ownership..."
    sudo chown caddy:caddy /etc/caddy/Caddyfile
    sudo chown -R app:app /opt/app
                  
    echo "📦 Renaming binaries on the host..."
    if sudo [ -f /opt/app/app ]; then
      sudo mv /opt/app/app /opt/app/app.old
    fi
    sudo mv /opt/app/app.new /opt/app/app

    echo "🔄 Reloading systemd and 🛠️ Enabling services..."
    sudo systemctl daemon-reexec
    sudo systemctl daemon-reload
    sudo systemctl enable app
    sudo systemctl enable --now db-backup.timer
    sudo systemctl enable caddy

    echo "🚀 Restarting legacy-node and waiting for it to be active..."
    sudo systemctl restart legacy-node
    sudo systemctl is-active --quiet legacy-node || (sudo journalctl -u legacy-node --no-pager -n 50 && exit 1)

    echo "✅ legacy-node is running. Starting app..."
    sudo systemctl restart app
    sudo systemctl restart caddy

    # Remove the old binary on the host if it exists
    if sudo [ -f /opt/app/app.old ]; then
      sudo rm /opt/app/app.old
    fi
EOF



echo "Deployment completed successfully."
