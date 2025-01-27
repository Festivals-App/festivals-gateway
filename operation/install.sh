#!/bin/bash
#
# install.sh 1.0.0
#
# Enables the firewall, installs the newest festivals-gateway and starts it as a service.
#
# (c)2020-2025 Simon Gaus
#

# Test for web server user
#
WEB_USER="www-data"
id -u "$WEB_USER" &>/dev/null;
if [ $? -ne 0 ]; then
  WEB_USER="www"
  if [ $? -ne 0 ]; then
    echo "Failed to find user to run web server. Exiting."
    exit 1
  fi
fi

# Move to working dir
#
mkdir -p /usr/local/festivals-gateway/install || { echo "Failed to create working directory. Exiting." ; exit 1; }
cd /usr/local/festivals-gateway/install || { echo "Failed to access working directory. Exiting." ; exit 1; }
echo "Installing festivals-gateway using port 8080."
sleep 1

# Get system os
#
if [ "$(uname -s)" = "Darwin" ]; then
  os="darwin"
elif [ "$(uname -s)" = "Linux" ]; then
  os="linux"
else
  echo "System is not Darwin or Linux. Exiting."
  exit 1
fi

# Get systems cpu architecture
#
if [ "$(uname -m)" = "x86_64" ]; then
  arch="amd64"
elif [ "$(uname -m)" = "arm64" ]; then
  arch="arm64"
else
  echo "System is not x86_64 or arm64. Exiting."
  exit 1
fi

# Build url to latest binary for the given system
#
file_url="https://github.com/Festivals-App/festivals-gateway/releases/latest/download/festivals-gateway-$os-$arch.tar.gz"
echo "The system is $os on $arch."
sleep 1

# Install festivals-gateway to /usr/local/bin/festivals-gateway. TODO: Maybe just link to /usr/local/bin?
#
echo "Downloading newest festivals-gateway binary release..."
curl -L "$file_url" -o festivals-gateway.tar.gz
tar -xf festivals-gateway.tar.gz
mv festivals-gateway /usr/local/bin/festivals-gateway || { echo "Failed to install festivals-gateway binary. Exiting." ; exit 1; }
echo "Installed the festivals-gateway binary to '/usr/local/bin/festivals-gateway'."
sleep 1

## Install server config file
mv config_template.toml /etc/festivals-gateway.conf
echo "Moved default festivals-gateway config to '/etc/festivals-gateway.conf'."
sleep 1

## Prepare log directory
mkdir /var/log/festivals-gateway || { echo "Failed to create log directory. Exiting." ; exit 1; }
echo "Create log directory at '/var/log/festivals-gateway'."

## Prepare server update workflow
mv update.sh /usr/local/festivals-gateway/update.sh
chmod +x /usr/local/festivals-gateway/update.sh
cp /etc/sudoers /tmp/sudoers.bak
echo "$WEB_USER ALL = (ALL) NOPASSWD: /usr/local/festivals-gateway/update.sh" >> /tmp/sudoers.bak
# Check syntax of the backup file to make sure it is correct.
visudo -cf /tmp/sudoers.bak
if [ $? -eq 0 ]; then
  # Replace the sudoers file with the new only if syntax is correct.
  sudo cp /tmp/sudoers.bak /etc/sudoers
else
  echo "Could not modify /etc/sudoers file. Please do this manually." ; exit 1;
fi

# Enable and configure the firewall.
#
if command -v ufw > /dev/null; then

  mv ufw_app_profile /etc/ufw/applications.d/festivals-gateway
  ufw allow festivals-gateway >/dev/null
  echo "Added festivals-gateway to ufw using port 443."
  sleep 1

elif ! [ "$(uname -s)" = "Darwin" ]; then
  echo "No firewall detected and not on macOS. Exiting."
  exit 1
fi

# Install systemd service
#
if command -v service > /dev/null; then

  if ! [ -f "/etc/systemd/system/festivals-gateway.service" ]; then
    mv service_template.service /etc/systemd/system/festivals-gateway.service
    echo "Created systemd service."
    sleep 1
  fi

  systemctl enable festivals-gateway > /dev/null
  echo "Enabled systemd service."
  sleep 1

elif ! [ "$(uname -s)" = "Darwin" ]; then
  echo "Systemd is missing and not on macOS. Exiting."
  exit 1
fi

## Set appropriate permissions
##
chown -R "$WEB_USER":"$WEB_USER" /usr/local/festivals-gateway
chown -R "$WEB_USER":"$WEB_USER" /var/log/festivals-gateway
chown "$WEB_USER":"$WEB_USER" /etc/festivals-gateway.conf

# Download FestivalsApp Root CA certificate
#--> to /usr/local/festivals-gateway/ca.crt

# Remving unused files
#
echo "Cleanup..."
cd /usr/local/festivals-gateway || exit
rm -R /usr/local/festivals-gateway/install
sleep 1

echo "Done!"
sleep 1

echo "You can start the server manually by running 'sudo systemctl start festivals-gateway' after you updated the configuration file at '/etc/festivals-gateway.conf'"
sleep 1
