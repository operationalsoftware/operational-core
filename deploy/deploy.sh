#!/bin/bash

# Exit if any command fails
set -e

# Check if environment is provided
if [ -z "$1" ]; then
  echo "Error: No deployment environment specified."
  echo "Usage: $0 <environment>"
  exit 1
fi

deployment_env="$1"

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

# Load in variables:
#   DEPLOY_HOST
#   DEPLOY_SSH_KEY_FLAG
#   DEPLOY_REQUIRES_CONFIRMATION
source ./read_config.sh "$deployment_env"

echo "Starting deployment to $DEPLOY_HOST"

if [ "$DEPLOY_REQUIRES_CONFIRMATION" = "true" ]; then
  echo "⚠️  You are about to deploy to: $deployment_env"
  read -p "Type the environment name to confirm: " confirm
  if [ "$confirm" != "$deployment_env" ]; then
    echo "Deployment cancelled."
    exit 1
  fi
fi

# Get the directory of this script
deploy_dir="$(cd "$(dirname "$0")" && pwd)"

# Set the working directory to that of this script
cd "$deploy_dir"


# Go to the main project directory
cd "$deploy_dir/.."

# Build the Golang binary
./build.sh app linux/amd64

# Copy built app to server
scp $DEPLOY_SSH_KEY_FLAG ./app "$DEPLOY_HOST:~"
scp $DEPLOY_SSH_KEY_FLAG "./deploy/$deployment_env.env" "$DEPLOY_HOST:~/.env"

# Remove the local binary
rm ./app

# Move back into this directory
cd "$deploy_dir"

# Copy config files to server
scp $DEPLOY_SSH_KEY_FLAG ./Caddyfile "$DEPLOY_HOST:~"
scp $DEPLOY_SSH_KEY_FLAG ./caddy.service "$DEPLOY_HOST:~"
scp $DEPLOY_SSH_KEY_FLAG ./app.service "$DEPLOY_HOST:~"
scp $DEPLOY_SSH_KEY_FLAG ./db-backup.service "$DEPLOY_HOST:~"
scp $DEPLOY_SSH_KEY_FLAG ./db-backup.sh "$DEPLOY_HOST:~"
scp $DEPLOY_SSH_KEY_FLAG ./db-backup.timer "$DEPLOY_HOST:~"
scp $DEPLOY_SSH_KEY_FLAG ./rclone.conf "$DEPLOY_HOST:~"

# Running commands in remote server via ssh
ssh $DEPLOY_SSH_KEY_FLAG "$DEPLOY_HOST" <<EOF
    set -e

    echo "📦 Ensuring directories exist..."
    sudo mkdir -p /opt/app
    sudo mkdir -p /home/app/.config/rclone

    echo "📦 Moving config and app files..."
    sudo mv ./Caddyfile /etc/caddy/Caddyfile
    sudo mv ./caddy.service /etc/systemd/system/caddy.service
    sudo mv ./app.service /etc/systemd/system/app.service
    sudo mv ./db-backup.service /etc/systemd/system/db-backup.service
    sudo mv ./db-backup.timer /etc/systemd/system/db-backup.timer
    sudo mv ./rclone.conf /home/app/.config/rclone/rclone.conf
    sudo mv ./app /opt/app/app.new
    sudo mv ./.env /opt/app/.env
    sudo mv ./db-backup.sh /opt/app/db-backup.sh

    echo "📦 Setting ownership..."
    sudo chown caddy:caddy /etc/caddy/Caddyfile
    sudo chown -R app:app /opt/app
    sudo chown -R app:app /home/app/.config/rclone

    echo "📦 Changing file permissions..."
    sudo chmod +x /opt/app/db-backup.sh
                  
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

    echo "✅ Starting app..."
    sudo systemctl restart app
    sudo systemctl restart caddy

    # Remove the old binary on the host if it exists
    if sudo [ -f /opt/app/app.old ]; then
      sudo rm /opt/app/app.old
    fi
EOF



echo "Deployment completed successfully."
