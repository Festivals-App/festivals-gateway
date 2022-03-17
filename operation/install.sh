#!/bin/bash
#
# install.sh 1.0.0
#
# Enables the firewall, installs the newest go and the festivals-gateway and starts it as a service.
#
# (c)2020-2022 Simon Gaus
#

# Move to working dir
#
mkdir /usr/local/festivals-gateway || { echo "Failed to create working directory. Exiting." ; exit 1; }
cd /usr/local/festivals-gateway || { echo "Failed to access working directory. Exiting." ; exit 1; }

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
mv festivals-server /usr/local/bin/festivals-gateway || { echo "Failed to install festivals-gateway binary. Exiting." ; exit 1; }
echo "Installed the festivals-gateway binary to '/usr/local/bin/festivals-gateway'."
mv config_template.toml /etc/festivals-gateway.conf
echo "Moved default festivals-gateway config to '/etc/festivals-gateway.conf'."
sleep 1

# Enable and configure the firewall.
#
if command -v ufw > /dev/null; then

  ufw allow 8080/tcp >/dev/null
  echo "Added festivals-gateway to ufw using port 8080."
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

# Remving unused files
#
echo "Cleanup..."
cd /usr/local || exit
rm -R /usr/local/festivals-gateway
sleep 1

echo "Done!"
sleep 1

echo "You can start the server manually by running 'systemctl start festivals-gateway' after you updated the configuration file at '/etc/festivals-gateway.conf'"
sleep 1
