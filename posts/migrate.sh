#!/bin/bash

# Load DB credentials from .env file
source .env

# posts always targets its own logical database (flashat_posts) rather
# than a DB_NAME from .env — that name is a fixed constant, not
# per-environment config. See config.PostsDBName.
migrate -path db/migrations \
  -database "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/flashat_posts?sslmode=disable" \
  "$@"
