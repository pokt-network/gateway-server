#!/bin/bash

# Set working directory to one level behind
cd "$(dirname "$0")/.."

# Parse command line arguments
if [[ $# -eq 0 ]]; then
  echo "Error: Missing required argument --name, --up, or --down"
  exit 1
fi

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
    shift # past argument
    ;;
    -d|--down)
    DOWN="true"
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
  # Migrate up
  migrate -database "$DB_CONNECTION_URL" -path "db_migrations" up 1
elif [ "$DOWN" = "true" ]; then
  # Migrate down
  migrate -database "$DB_CONNECTION_URL" -path "db_migrations" down 1
else
  # Create new migration
  if [ -z "$NAME" ]; then
    echo "Error: Missing required argument --name"
    exit 1
  fi
  migrate -database "$DB_CONNECTION_URL" create -ext sql -seq -dir "db_migrations" "$NAME"
fi