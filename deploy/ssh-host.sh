#!/bin/bash

# Exit if any command fails
set -e

# Get the directory of this script
this_dir="$(cd "$(dirname "$0")" && pwd)"

# Set the working directory to that of this script
cd $this_dir

# Source the get-deploy-config script to read the configuration values
source ./get-deploy-config.sh

# Execute the SSH command with the optional SSH key
ssh $ssh_key_flag "$host"

