#!/bin/bash

error_exit() {
    echo "Error: $1"
    exit 1
}

check_success() {
    if [ $? -ne 0 ]; then
        error_exit "$1"
    else
        echo "Success: $2"
    fi
}

echo "=== Virtualization Server Setup ==="

echo -n "Adding community repository... "
echo "http://dl-cdn.alpinelinux.org/alpine/v3.21/community" >> /etc/apk/repositories
check_success "Failed to add repository" "Repository added"

echo -n "Installing required packages... "
cd /root/VirtualizationServer || error_exit "Project directory not found"
[ -f installed_packages.txt ] || error_exit "installed_packages.txt not found"
apk update
cat installed_packages.txt | xargs apk add
check_success "Package installation failed" "Packages installed"

echo -n "Installing TUNEL client... "
curl -sSLf https://get.tuna.am | sh
check_success "TUNEL installation failed" "TUNEL installed"

read -p "Please write your tuna token: " tuna_token
echo -n "Saving TUNEL token... "
tuna config save-token "$tuna_token"
check_success "Failed to save token" "Token saved"

read -p "Please write your domain: " domain_name

TARGET_SCRIPT="/root/VirtualizationServer/vm-api/scripts/start_tuna.sh"
echo -n "Updating domain in scripts... "
sed -i "s/<yourdomain>/$domain_name/g" "$TARGET_SCRIPT"
check_success "Domain update failed" "Domain updated"

echo -n "Setting script permissions... "
chmod +x vm-api/scripts/*
check_success "Permission setting failed" "Permissions set"

echo -n "Building application... "
cd vm-api/ && go build -o run_to_start
check_success "Build failed" "Application built"

echo -n "Configuring cron jobs... "
(crontab -l 2>/dev/null; echo "@reboot sleep 5 && cd /root/VirtualizationServer/vm-api/scripts/ && ./start_tuna.sh") | crontab -
(crontab -l 2>/dev/null; echo "0 */4 * * * cd /root/VirtualizationServer/vm-api/scripts/ && ./check_qemu.sh") | crontab -
(crontab -l 2>/dev/null; echo "0 */4 * * * cd /root/VirtualizationServer/vm-api/scripts/ && ./clear_ports.sh") | crontab -
check_success "Cron configuration failed" "Cron jobs configured"
mkdir /root/VirtualizationServer/image/vms

echo "===================================="
echo "Setup completed successfully!"
echo "To start work please reboot system and after reboot run file:"
echo "/root/VirtualizationServer/vm-api/run_to_start"
echo "===================================="

exit 0
