#!/bin/bash

# Load environment variables
set -a
source .env
set +a

# Run migration
echo "Running database migrations..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f internal/infrastructure/persistence/migrations/001_create_users_table.sql

if [ $? -eq 0 ]; then
    echo "Migration completed successfully!"
else
    echo "Migration failed!"
    exit 1
fi
