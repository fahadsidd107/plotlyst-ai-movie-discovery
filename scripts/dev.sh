#!/bin/bash

set -e

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SERVER_DIR="$PROJECT_ROOT/Server/PlotlystServer"
CLIENT_DIR="$PROJECT_ROOT/Client/plotlyst-client"
export PATH="/opt/homebrew/bin:/opt/homebrew/opt/node@24/bin:$PATH"

for command_name in go node npm mongosh mongoimport brew; do
  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "Missing required command: $command_name"
    exit 1
  fi
done

if ! mongosh --quiet --nodb --eval 'quit(0)' >/dev/null 2>&1; then
  echo "MongoDB shell is not working."
  exit 1
fi

if ! mongosh --quiet "mongodb://127.0.0.1:27017/admin" --eval 'db.runCommand({ping: 1}).ok' >/dev/null 2>&1; then
  echo "Starting MongoDB..."
  brew services start mongodb-community@8.0 >/dev/null
  for attempt in 1 2 3 4 5 6 7 8 9 10; do
    if mongosh --quiet "mongodb://127.0.0.1:27017/admin" --eval 'db.runCommand({ping: 1}).ok' >/dev/null 2>&1; then
      break
    fi
    sleep 1
  done
fi

seed_collection() {
  collection_name="$1"
  source_file="$2"
  document_count="$(mongosh --quiet "mongodb://127.0.0.1:27017/plotlyst" --eval "db.getCollection('$collection_name').estimatedDocumentCount()")"
  if [ "$document_count" = "0" ]; then
    echo "Importing $collection_name sample data..."
    mongoimport --quiet --uri="mongodb://127.0.0.1:27017" --db=plotlyst --collection="$collection_name" --file="$PROJECT_ROOT/plotlyst-seed-data/$source_file" --jsonArray
  fi
}

seed_collection movies movies.json
seed_collection genres genres.json
seed_collection rankings rankings.json
seed_collection users users.json

if [ ! -d "$CLIENT_DIR/node_modules" ]; then
  echo "Installing frontend packages (first run only)..."
  npm --prefix "$CLIENT_DIR" ci
fi

cleanup() {
  trap - EXIT INT TERM
  [ -n "$SERVER_PID" ] && kill "$SERVER_PID" 2>/dev/null || true
  [ -n "$CLIENT_PID" ] && kill "$CLIENT_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

echo "Starting Plotlyst backend at http://localhost:8080"
(cd "$SERVER_DIR" && go run .) &
SERVER_PID=$!

echo "Starting Plotlyst frontend at http://localhost:5173"
npm --prefix "$CLIENT_DIR" run dev &
CLIENT_PID=$!

echo "Press Ctrl+C to stop both servers."
wait "$SERVER_PID" "$CLIENT_PID"
