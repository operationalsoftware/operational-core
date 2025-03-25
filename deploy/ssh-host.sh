#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Source the get-deploy-config script to read the configuration values
source ./get-deploy-config.sh

# Execute the SSH command with the optional SSH key
ssh $ssh_key_flag "$host"

