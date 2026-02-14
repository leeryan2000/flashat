#!/bin/bash

# 1. Load configuration
source .env
echo "Targeting deployment for $APP_NAME..."

# Define functions for reuse
deploy_frontend() {
    echo "--- Starting Frontend Deployment ---"
    cd web || exit
    npm run build
    if [ $? -eq 0 ]; then
        echo "Uploading frontend dist to $FRONTEND_DEST..."
        scp -i ../$PEM_FILE -r dist/* $SERVER_USER@$SERVER_IP:$FRONTEND_DEST
    else
        echo "Frontend build failed!"
    fi
    cd ..
}

deploy_backend() {
    echo "--- Starting Backend Deployment ---"
    echo "Compiling Go binary for Linux..."
    GOOS=linux GOARCH=amd64 go build -o bin/$APP_NAME main.go
    if [ $? -eq 0 ]; then
        echo "Uploading binary to $BACKEND_DEST..."
        scp -i $PEM_FILE bin/$APP_NAME $SERVER_USER@$SERVER_IP:$BACKEND_DEST
    else
        echo "Backend build failed!"
    fi
}

# 2. Logic to handle arguments
# $1 refers to the first word after the script name
case "$1" in
    frontend)
        deploy_frontend
        ;;
    backend)
        deploy_backend
        ;;
    all | "")
        deploy_frontend
        deploy_backend
        ;;
    *)
        echo "Usage: $0 {all|frontend|backend}"
        ;;
esac

echo "Deployment task finished."
read -p "Press [Enter] key to exit..."