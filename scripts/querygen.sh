# Set working directory to one level behind
cd "$(dirname "$0")/.."

# Load environment variables from .env file located one level behind
if [ -f .env ]; then
  export $(cat .env | sed 's/#.*//g' | xargs)
else
  echo "Error: .env file not found in parent directory."
  exit 1
fi

pggen gen go \
    --postgres-connection "$DB_CONNECTION_URL" \
    --query-glob internal/db_query/*.sql \
    --go-type 'int8=int' \
    --go-type 'text=string'