#!/bin/bash

export PGPASSWORD="dev"

# Define your database and user
DB_NAME="dev"
DB_USER="dev"

psql -U "$DB_USER" -d "$DB_NAME" -f db/cleanup.sql

echo "✅ Database cleaned up."