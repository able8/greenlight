#!/bin/bash
set -eu

TIMEZONE="Asia/Shanghai"
USERNAME="greenlight"

# Prompt to enter a password for the PostgreSQL greenlight user
read -p "Enter password for greenlight DB user: " DB_PASSWORD

# Force all output to be presented in en_US for the duration of this script.
# This avoids any "setting locale failed" errors while this script is running,
# before we have installed support for all locale.
export LC_ALL=en_US.UTF-8


# Enable the "universe" repository
# This repository consists free and open source software but Ubuntu doesn’t guarantee of regular security updates to software in this category.
add-apt-repository --yes universe
apt update

# Update all software packages. Using the --force-conf
apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

# Set the system timezone and install all locales.
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

# Add the new user (and give them sudo privileges).
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"

# Force a password to be set for the new user the first time they log in
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

# Copy the SSH keys from the root user to the new User.
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}

# Configure the firewall to allow ssh, http and https traffic.
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# Install fail2ban
apt --yes install fail2ban

# Install the migrate CLI tool.
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 /usr/local/bin/migrate

# Install PostgreSQL
apt --yes install postgresql

# Set up the greenlight DB and create a user account with the password entered earlier.
sudo -i -u postgresql psql -c "CREATE DATABASE greenlight"
sudo -i -u postgresql psql -d greenlight -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgresql psql -d greenlight -c "CREATE ROLE greenlight WITH LOGIN PASSWORD '${DB_PASSWORD}'"

# Add a DSN for connection to the greenlight database to the system-wide environment variables in the /etc/environment file.
echo "export GREENLIGHT_DB_DSN='postgres://greenlight:${DB_PASSWORD}@localhost:/greenlight?sslmode=disable'" >> /etc/environment


apt update
apt --yes install caddy

echo "Script complete! Rebooting..."
# reboot
