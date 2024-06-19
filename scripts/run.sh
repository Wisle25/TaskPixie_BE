#!/bin/sh

# Load environment variables from .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Check if APP_ENV is set to prod or dev
if [ "$APP_ENV" = "prod" ]; then
  make migrate-up
fi

# Run the main application
./main
