#!/bin/bash

set -e;

timestamp=$(date +%F_%T | tr ':' '-');
echo "Starting backup at $timestamp";

dump_file="/opt/app/${DUMP_PREFIX}-${timestamp}.dump"
orbit_object="${DUMP_PREFIX}-${timestamp}.dump"

# Create database dump
pg_dump -Fc --no-acl -U postgres $PG_DATABASE > "$dump_file"

# Upload to Orbit
rclone copyto "$dump_file" "orbit:$ORBIT_BACKUP_CONTAINER/$orbit_object"

timestamp=$(date +%F_%T | tr ':' '-');
echo "Backup complete at $timestamp - uploaded to orbit:$ORBIT_BACKUP_CONTAINER/$orbit_object"

rm "$dump_file"

rclone delete "orbit:$ORBIT_BACKUP_CONTAINER" --min-age 30d
