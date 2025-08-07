#!/bin/bash

# Load DB credentials from .env file
source .env

# Run migrate command with any arguments passed after ./migrate.sh
migrate -path db/migrations \
  -database "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
  "$@"