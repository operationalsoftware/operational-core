#!/bin/bash

set -euo pipefail

echo "----- Updating System -----"
sudo apt update && sudo apt upgrade -y

echo "----- Installing Required Packages -----"
sudo apt install -y curl git rsync gnupg2 net-tools ca-certificates lsb-release apt-transport-https

echo "----- Installing Chromium -----"
sudo apt install -y chromium chromium-driver

echo "----- Installing Node.js (LTS) -----"
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

echo "----- Installing Rclone -----"
curl https://downloads.rclone.org/v1.69.2/rclone-v1.69.2-linux-amd64.deb -O
sudo dpkg -i rclone-v1.69.2-linux-amd64.deb
sudo rm rclone-v1.69.2-linux-amd64.deb

echo "----- Installing Caddy -----"
sudo apt install -y debian-keyring debian-archive-keyring
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | \
  sudo tee /etc/apt/trusted.gpg.d/caddy.gpg > /dev/null
echo "deb [trusted=yes] https://dl.cloudsmith.io/public/caddy/stable/deb/debian any-version main" | \
  sudo tee /etc/apt/sources.list.d/caddy.list
sudo apt update
sudo apt install -y caddy

echo "----- Creating 'app' System User -----"
sudo useradd --system --create-home --shell /bin/bash app
echo 'app:app' | sudo chpasswd
echo 'app ALL=(ALL) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/app
sudo chmod 440 /etc/sudoers.d/app

echo "----- Creating '/home/app/.pgpass' file -----"
cat <<EOF > /home/app/.pgpass
localhost:5432:*:postgres:postgres
EOF
sudo chown app:app /home/app/.pgpass
sudo chmod 600 /home/app/.pgpass

echo "----- Creating Rclone config file -----"
mkdir -p /home/app/.config/rclone
# UPDATE following config from orbit container
cat <<EOF > /home/app/.config/rclone/rclone.conf
[orbit]
type = swift
tenant = acc-abcdef
user = cli-abcdef
key = abcdefghi
auth = https://orbit.brightbox.com/v3
EOF
sudo chown app:app /home/app/.config/rclone/rclone.conf

echo "----- Installing PostgreSQL 16 -----"
echo "----- Adding PostgreSQL APT Repository -----"
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

sudo apt update
sudo apt install -y postgresql-16

echo "----- Enabling and Starting PostgreSQL -----"
sudo systemctl enable postgresql
sudo systemctl start postgresql

echo "----- Setting Password for 'postgres' User -----"
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"

APP_DB_NAME="batten_allen"

echo "----- Creating PostgreSQL User and Database -----"
sudo -u postgres psql <<EOF
CREATE DATABASE "$APP_DB_NAME" OWNER postgres;
EOF

# Ensure only local access
PG_HBA="/etc/postgresql/16/main/pg_hba.conf"
POSTGRES_CONF="/etc/postgresql/16/main/postgresql.conf"

echo "----- Configuring PostgreSQL for Local-Only Access -----"
sudo sed -i "s/^#listen_addresses = 'localhost'/listen_addresses = 'localhost'/" "$POSTGRES_CONF"
sudo sed -i "s/^listen_addresses = '\*'/listen_addresses = 'localhost'/" "$POSTGRES_CONF"

# Change local authentication from peer to md5 for password login (optional)
sudo sed -i "s/^local\s\+all\s\+all\s\+peer/local all all md5/" "$PG_HBA"
sudo sed -i "s/^local\s\+all\s\+postgres\s\+peer/local all postgres md5/" "$PG_HBA"

# Enable logging collector (required for logging to file)
sudo sed -i "s/^#logging_collector = off/logging_collector = on/" "$POSTGRES_CONF"

# Set log directory and file name (optional but recommended)
sudo sed -i "s|^#log_directory = 'log'|log_directory = 'log'|" "$POSTGRES_CONF"
sudo sed -i "s/^#log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'/log_filename = 'postgresql-%Y-%m-%d.log'/" "$POSTGRES_CONF"

# Set log_min_duration_statement to 500ms
sudo sed -i "s/^#log_min_duration_statement = -1/log_min_duration_statement = 500/" "$POSTGRES_CONF"
sudo sed -i "s/^log_min_duration_statement = -1/log_min_duration_statement = 500/" "$POSTGRES_CONF"


sudo systemctl restart postgresql

echo "----- Setup Complete! -----"
echo "Node version: $(node -v)"
echo "NPM version: $(npm -v)"
echo "Caddy version: $(caddy version)"
echo "Rclone version: $(rclone --version)"