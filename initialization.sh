#!/bin/bash

if [ -z "$1" ]; then
  echo "Error: Please write your domain"
  echo "Example: ./initialization.sh your_domain"
  exit 1
fi

TARGET_SCRIPT="/root/VirtualizationServer/vm-api/scripts/start_tuna.sh"


sed -i "s/<yourdomain>/$1/g" "$TARGET_SCRIPT"

if grep -q "$1" "$TARGET_SCRIPT"; then
  echo "Domain changed to: $1"
else
  echo "Something went wrong!"
  exit 1
fi
