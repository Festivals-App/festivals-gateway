#!/bin/bash
#
# update.sh 1.0.0
#
# Updates the festivals-gateway and restarts it.
#
# (c)2020-2022 Simon Gaus
#

# Move to working dir
#
mkdir /usr/local/festivals-gateway/install || { echo "Failed to create working directory. Exiting." ; exit 1; }
cd /usr/local/festivals-gateway/install || { echo "Failed to access working directory. Exiting." ; exit 1; }

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

# Updating festivals-gateway to the newest binary release
#
echo "Downloading newest festivals-gateway binary release..."
curl -L "$file_url" -o festivals-gateway.tar.gz
tar -xf festivals-gateway.tar.gz
mv festivals-gateway /usr/local/bin/festivals-gateway || { echo "Failed to install festivals-gateway binary. Exiting." ; exit 1; }
echo "Updated festivals-gateway binary."
sleep 1

# Removing unused files
#
echo "Cleanup..."
cd /usr/local/festivals-gateway || { echo "Failed to access server directory. Exiting." ; exit 1; }
rm -r /usr/local/festivals-gateway/install
sleep 1

# Restart the festivals-gateway
#
systemctl restart festivals-gateway
echo "Restarted the festivals-gateway"
sleep 1

echo "Done!"
sleep 1 