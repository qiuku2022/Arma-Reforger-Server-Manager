#!/bin/bash

# ARSM Build Script
# 构建 Arma Reforger Server Manager

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "=== Building ARSM ==="

# 构建前端
echo "[1/3] Building frontend..."
cd frontend
npm install
npm run build
cd ..

# 复制前端产物到后端
echo "[2/3] Copying static files..."
rm -rf backend/static/*
cp -r frontend/dist/* backend/static/

# 构建后端
echo "[3/3] Building backend..."
cd backend

# 检测目标平台
if [[ "$1" == "windows" ]]; then
    echo "Building for Windows..."
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o arsm.exe .
    echo "Output: backend/arsm.exe"
elif [[ "$1" == "all" ]]; then
    echo "Building for all platforms..."
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o arsm-linux-amd64 .
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o arsm-windows-amd64.exe .
    echo "Output: backend/arsm-linux-amd64, backend/arsm-windows-amd64.exe"
else
    echo "Building for current platform..."
    go build -ldflags="-s -w" -o arsm .
    echo "Output: backend/arsm"
fi

cd ..
echo "=== Build complete ==="
