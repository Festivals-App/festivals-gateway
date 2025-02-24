#!/bin/bash
#
# install.sh - FestivalsApp Gateway Installer Script
#
# Enables the firewall, installs the newest FestivalsApp Gateway and starts it as a service.
#
# (c)2020-2025 Simon Gaus
#

# ─────────────────────────────────────────────────────────────────────────────
# 🔍 Detect Web Server User
# ─────────────────────────────────────────────────────────────────────────────
WEB_USER="www-data"
if ! id -u "$WEB_USER" &>/dev/null; then
    WEB_USER="www"
    if ! id -u "$WEB_USER" &>/dev/null; then
        echo -e "\n\033[1;31m❌  ERROR: Web server user not found! Exiting.\033[0m\n"
        exit 1
    fi
fi

# ─────────────────────────────────────────────────────────────────────────────
# 📁 Setup Working Directory
# ─────────────────────────────────────────────────────────────────────────────
WORK_DIR="/usr/local/festivals-gateway/install"
mkdir -p "$WORK_DIR" && cd "$WORK_DIR" || { echo -e "\n\033[1;31m❌  ERROR: Failed to create/access working directory!\033[0m\n"; exit 1; }
echo -e "\n📂  Working directory set to \e[1;34m$WORK_DIR\e[0m"
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🖥  Detect System OS and Architecture
# ─────────────────────────────────────────────────────────────────────────────
if [ "$(uname -s)" = "Darwin" ]; then
    os="darwin"
elif [ "$(uname -s)" = "Linux" ]; then
    os="linux"
else
    echo -e "\n🚨  ERROR: Unsupported OS. Exiting.\n"
    exit 1
fi
if [ "$(uname -m)" = "x86_64" ]; then
    arch="amd64"
elif [ "$(uname -m)" = "arm64" ]; then
    arch="arm64"
else
    echo -e "\n🚨  ERROR: Unsupported CPU architecture. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# 📦 Install FestivalsApp Gateway
# ─────────────────────────────────────────────────────────────────────────────
file_url="https://github.com/Festivals-App/festivals-gateway/releases/latest/download/festivals-gateway-$os-$arch.tar.gz"
echo -e "\n📥  Downloading latest FestivalsApp Gateway release..."
curl --progress-bar -L "$file_url" -o festivals-gateway.tar.gz
echo -e "📦  Extracting binary..."
tar -xf festivals-gateway.tar.gz
mv festivals-gateway /usr/local/bin/festivals-gateway || {
    echo -e "\n🚨  ERROR: Failed to install FestivalsApp Gateway binary. Exiting.\n"
    exit 1
}
echo -e "✅  Installed FestivalsApp Gateway to \e[1;34m/usr/local/bin/festivals-gateway\e[0m."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🛠  Install Server Configuration File
# ─────────────────────────────────────────────────────────────────────────────
mv config_template.toml /etc/festivals-gateway.conf
if [ -f "/etc/festivals-gateway.conf" ]; then
    echo -e "✅  Configuration file moved to \e[1;34m/etc/festivals-gateway.conf\e[0m."
else
    echo -e "\n🚨  ERROR: Failed to move configuration file. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 📂  Prepare Log Directory
# ─────────────────────────────────────────────────────────────────────────────
mkdir -p /var/log/festivals-gateway || {
    echo -e "\n🚨  ERROR: Failed to create log directory. Exiting.\n"
    exit 1
}
echo -e "✅  Log directory created at \e[1;34m/var/log/festivals-gateway\e[0m."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔄 Prepare Remote Update Workflow
# ─────────────────────────────────────────────────────────────────────────────
mv update.sh /usr/local/festivals-gateway/update.sh
chmod +x /usr/local/festivals-gateway/update.sh
cp /etc/sudoers /tmp/sudoers.bak
echo "$WEB_USER ALL = (ALL) NOPASSWD: /usr/local/festivals-gateway/update.sh" >> /tmp/sudoers.bak
# Validate and replace sudoers file if syntax is correct
if visudo -cf /tmp/sudoers.bak &>/dev/null; then
    sudo cp /tmp/sudoers.bak /etc/sudoers
    echo -e "✅  Prepared remote update workflow."
else
    echo -e "\n🚨  ERROR: Could not modify /etc/sudoers file. Please do this manually. Exiting.\n"
    exit 1
fi
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🔥 Enable and Configure Firewall
# ─────────────────────────────────────────────────────────────────────────────
if command -v ufw > /dev/null; then
    echo -e "\n🚀  Configuring UFW firewall..."
    mv ufw_app_profile /etc/ufw/applications.d/festivals-gateway
    ufw allow festivals-gateway >/dev/null
    echo -e "✅  Added festivals-gateway to UFW with port 443."
    sleep 1
elif ! [ "$(uname -s)" = "Darwin" ]; then
    echo -e "\n🚨  ERROR: No firewall detected and not on macOS. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# ⚙️  Install Systemd Service
# ─────────────────────────────────────────────────────────────────────────────
if command -v service > /dev/null; then
    echo -e "\n🚀  Configuring systemd service..."
    if ! [ -f "/etc/systemd/system/festivals-gateway.service" ]; then
        mv service_template.service /etc/systemd/system/festivals-gateway.service
        echo -e "✅  Created systemd service configuration."
        sleep 1
    fi
    systemctl enable festivals-gateway > /dev/null
    echo -e "✅  Enabled systemd service for FestivalsApp Gateway."
    sleep 1
elif ! [ "$(uname -s)" = "Darwin" ]; then
    echo -e "\n🚨  ERROR: Systemd is missing and not on macOS. Exiting.\n"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# 🔑 Set Appropriate Permissions
# ─────────────────────────────────────────────────────────────────────────────
chown -R "$WEB_USER":"$WEB_USER" /usr/local/festivals-gateway
chown -R "$WEB_USER":"$WEB_USER" /var/log/festivals-gateway
chown "$WEB_USER":"$WEB_USER" /etc/festivals-gateway.conf
echo -e "\n🔐  Set Appropriate Permissions."
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🧹 Cleanup Installation Files
# ─────────────────────────────────────────────────────────────────────────────
echo -e "🧹  Cleaning up installation files..."
cd /usr/local/festivals-gateway || exit
rm -R /usr/local/festivals-gateway/install
sleep 1

# ─────────────────────────────────────────────────────────────────────────────
# 🎉 Final Message
# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"
echo -e "\033[1;32m✅  INSTALLATION COMPLETE! 🚀\033[0m"
echo -e "\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m"
echo -e "\n📂 \033[1;34mBefore starting, you need to:\033[0m"
echo -e "\n   \033[1;34m1. Configure the mTLS certificates.\033[0m"
echo -e "   \033[1;34m3. Configuring the FestivlasApp Root CA.\033[0m"
echo -e "   \033[1;34m4. Update the configuration file at:\033[0m"
echo -e "\n   \033[1;32m    /etc/festivals-gateway.conf\033[0m"
echo -e "\n🔹 \033[1;34mThen start the server manually:\033[0m"
echo -e "\n   \033[1;32m    sudo systemctl start festivals-gateway\033[0m"
echo -e "\n\033[1;32m══════════════════════════════════════════════════════════════════════════\033[0m\n"
sleep 1
