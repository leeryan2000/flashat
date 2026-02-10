#!/bin/bash

# 1. Set variables
APP_NAME="flashat-backend"
PEM_FILE="leeryan.pem"
SERVER_USER="ec2-user"
SERVER_IP="18.140.243.27"
DEST_PATH="/home/ec2-user/flashat/backend/"

echo "Starting deployment for $APP_NAME..."

# 2. Compile for Linux (The Go standard way)
echo "Building binary for Linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o $APP_NAME main.go

# 3. Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful. Transferring to EC2..."
    
    # 4. Transfer the binary
    scp -i $PEM_FILE $APP_NAME $SERVER_USER@$SERVER_IP:$DEST_PATH
    
    echo "Deployment complete!"
else
    echo "Build failed. Deployment aborted."
    exit 1
fi