#!/bin/bash
REPO_DIR=$(pwd)

# Обновление путей в скриптах
find . -type f -exec sed -i "s|/root/VirtualizationServer|$REPO_DIR|g" {} +

# Создание конфигурации для Go-кода
echo "VIRT_SERVER_ROOT=$REPO_DIR" > .env
