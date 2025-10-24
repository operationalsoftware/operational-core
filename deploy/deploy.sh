#!/bin/bash

# Exit if any command fails
set -e

# Allow override of git checks with SKIP_GIT_CHECKS=1
if [ -n "$SKIP_GIT_CHECKS" ]; then
  echo "⚠️  Git checks are being skipped (SKIP_GIT_CHECKS=$SKIP_GIT_CHECKS)"
fi

# Check if environment is provided
if [ -z "$1" ]; then
  echo "Error: No deployment environment specified."
  echo "Usage: $0 <environment>"
  exit 1
fi

deployment_env="$1"

# Check if there are uncommitted changes (unless overridden)
if [ -z "$SKIP_GIT_CHECKS" ] && [ -n "$(git status --porcelain)" ]; then
  echo "Error: There are uncommitted changes. Please commit or stash your changes before deploying."
  exit 1
fi

# Check if there are any unpushed commits (unless overridden)
if [ -z "$SKIP_GIT_CHECKS" ] && [ -n "$(git log origin/$(git rev-parse --abbrev-ref HEAD)..HEAD)" ]; then
  echo "Error: There are unpushed commits. Please push your changes before deploying."
  exit 1
fi

# Get the directory of this script
deploy_dir="$(cd "$(dirname "$0")" && pwd)"

# Set the working directory to that of this script
cd "$deploy_dir"

# Load in variables:
#   DEPLOY_HOST
#   DEPLOY_SSH_KEY_FLAG
#   DEPLOY_REQUIRES_CONFIRMATION
source ./read_config.sh "$deployment_env"

echo "Starting deployment to $DEPLOY_HOST"

if [ "$DEPLOY_REQUIRES_CONFIRMATION" = "true" ]; then
  echo "⚠️  You are about to deploy to: $deployment_env"
  read -p "Type \"$deployment_env\" to confirm: " confirm
  if [ "$confirm" != "$deployment_env" ]; then
    echo "Deployment cancelled."
    exit 1
  fi
fi

# Go to the main project directory
cd "$deploy_dir/.."

# Build the Golang binary
./build.sh app linux/amd64

# Copy built app to server
scp $DEPLOY_SSH_KEY_FLAG ./app "$DEPLOY_HOST:~"

# Prepare an environment snapshot from current shell (expected to be injected via Phase or other means)
# We only include variables used by the app for safety.
allowed_vars=(
  "APP_ENV" "GO_ENV" "SITE_ADDRESS"
  "PG_USER" "PG_PASSWORD" "PG_HOST" "PG_PORT" "PG_DATABASE"
  "MS_OAUTH_CLIENT_ID" "MS_OAUTH_SECRET" "MS_OAUTH_TENANT_ID"
  "SWIFT_API_USER" "SWIFT_API_KEY" "SWIFT_AUTH_URL" "SWIFT_TENANT_ID" "SWIFT_CONTAINER"
  "SECURE_COOKIE_HASH_KEY" "SECURE_COOKIE_BLOCK_KEY" "AES_256_ENCRYPTION_KEY"
  "SYSTEM_USER_PASSWORD"
  "DUMP_PREFIX" "ORBIT_BACKUP_CONTAINER"
)

tmp_env_file="/tmp/app.env"
rm -f "$tmp_env_file"
for key in "${allowed_vars[@]}"; do
  val=$(printenv "$key" || true)
  if [ -n "$val" ]; then
    # NOTE: values are written as-is; multiline values are not supported here
    printf '%s=%s\n' "$key" "$val" >> "$tmp_env_file"
  fi
done

# Upload environment snapshot to server home; the server step will move it into place
scp $DEPLOY_SSH_KEY_FLAG "$tmp_env_file" "$DEPLOY_HOST:~/app.env"
rm -f "$tmp_env_file"

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
    # Put environment in place for systemd units
    if [ -f ./app.env ]; then
      sudo mv ./app.env /opt/app/app.env
    fi
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
