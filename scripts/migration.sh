#!/bin/bash

# Set working directory to one level behind
cd "$(dirname "$0")/.."

# Parse command line arguments
if [[ $# -eq 0 ]]; then
  echo "Error: Missing required argument --name, --up, or --down"
  exit 1
fi

UP_MIGRATION_NUMBER="" # Default to applying all migrations for up
DOWN_MIGRATION_NUMBER="" # No default for down, must be explicitly provided

while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -n|--name)
    NAME="$2"
    shift # past argument
    shift # past value
    ;;
    -u|--up)
    UP="true"
    if [[ -n "$2" ]] && ! [[ "$2" =~ ^- ]]; then
        UP_MIGRATION_NUMBER="$2"
        shift
    fi
    shift # past argument
    ;;
    -d|--down)
    DOWN="true"
    if [[ -n "$2" ]]; then
        DOWN_MIGRATION_NUMBER="$2"
        shift
    else
        echo "Error: Down migration requires a specific number of migrations or -all flag to revert all migrations."
        exit 1
    fi
    shift # past argument
    ;;
    *)    # unknown option
    echo "Unknown option: $1"
    exit 1
    ;;
esac
done

# Load environment variables from .env file located one level behind
if [ -f .env ]; then
  export $(cat .env | sed 's/#.*//g' | xargs)
else
  echo "Error: .env file not found in parent directory."
  exit 1
fi

# Check if migrating up, down, or creating a new migration
if [ "$UP" = "true" ]; then
  # Migrate up to a number of steps or to the latest version
  migrate -database "$DB_CONNECTION_URL" -path "db_migrations" up ${UP_MIGRATION_NUMBER}
elif [ "$DOWN" = "true" ]; then
  # Migrate down to a number of steps or to the initial version
  migrate -database "$DB_CONNECTION_URL" -path "db_migrations" down ${DOWN_MIGRATION_NUMBER}
else
  # Create new migration
  if [ -z "$NAME" ]; then
    echo "Error: Missing required argument --name"
    exit 1
  fi
  migrate -database "$DB_CONNECTION_URL" create -ext sql -seq -dir "db_migrations" "$NAME"
fi
