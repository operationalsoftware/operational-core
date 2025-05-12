#!/bin/bash

set -e;

# Load Orbit credentials and config
source /opt/app/.env

TIMESTAMP=$(date +%F_%T | tr ':' '-');
echo "Starting backup at $TIMESTAMP";

DUMP_FILE="/opt/app/${DUMP_PREFIX}-${TIMESTAMP}.dump"
ORBIT_OBJECT="${DUMP_PREFIX}-${TIMESTAMP}.dump"

# Create database dump
pg_dump -Fc --no-acl -U postgres batten_allen > "$DUMP_FILE"

# Upload to Orbit
rclone copyto "$DUMP_FILE" "orbit:$ORBIT_BACKUP_CONTAINER/$ORBIT_OBJECT"

TIMESTAMP=$(date +%F_%T | tr ':' '-');
echo "Backup complete at $TIMESTAMP - uploaded to orbit:$ORBIT_BACKUP_CONTAINER/$ORBIT_OBJECT"

rm "$DUMP_FILE"

rclone delete "orbit:$ORBIT_BACKUP_CONTAINER" --min-age 90d