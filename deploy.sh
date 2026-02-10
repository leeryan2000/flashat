#!/bin/bash


# 1. Set variables
source .env
echo "Starting deployment for $APP_NAME..."

# 2. Compile for Linux
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
    read -p "Press [Enter] key to exit..."

    exit 1
fi

read -p "Press [Enter] key to exit..."