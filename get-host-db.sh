#!/bin/bash

# Usage: ./get-host-db.sh [local_backup_name]

# Exit immediately if a command exits with a non-zero status
set -e

# Validate the local backup name parameter
if [ -z "$1" ]; then
  echo "Error: Local backup name is required."
  exit 1
fi

local_backup_name="$1"

# Source the get-deploy-config script to read the configuration values
source ./get-deploy-config.sh

# Step 1: Define the remote database path and backup path
remote_db_path="/app/app.db"
remote_backup_path="/app/app_backup.db"

# Step 2: Create a backup on the remote server
echo "Creating a backup of the remote database..."
remote_command="sqlite3 '$remote_db_path' '.backup \"$remote_backup_path\"'"
ssh $ssh_key_flag "$host" "$remote_command"

# Step 3: Copy the remote backup to the local directory
echo "Copying the backup from the remote host to the local directory..."
scp $ssh_key_flag "$host:$remote_backup_path" "$local_backup_name"

# Step 4: Remove the remote backup file after copying
echo "Cleaning up the remote backup file..."
remote_cleanup_command="rm '$remote_backup_path'"
ssh $ssh_key_flag "$host" "$remote_cleanup_command"

echo "Database backup and copy process completed successfully."

