#!/bin/bash

set -e

echo "🔵 Starting LIST migration tests..."

function check_metana_binary() {
  if ! command -v metana &> /dev/null; then
    echo "❌ metana binary not found in PATH. Please install it first with: go install ./..."
    exit 1
  fi
  echo "✅ metana binary found at $(which metana)"
}

function create_temp_dir() {
  TEMP_DIR=$(mktemp -d)
  echo "✅ Created temp directory: $TEMP_DIR"
}

function run_metana_init() {
  echo "🚀 Running metana init..."
  metana init --dir "$TEMP_DIR/migrations"
}

function create_migrations() {
  echo "🚀 Creating migration files..."

  (cd "$TEMP_DIR" && metana create initSchema)
  sleep 1
  (cd "$TEMP_DIR" && metana create initSchema2)

  echo "✅ Created two migration scripts."
}

function run_metana_list_and_validate_empty_executed() {
  echo "📋 Running metana list (before running migrations)..."
  
  LIST_OUTPUT=$(cd "$TEMP_DIR" && metana list)

  if echo "$LIST_OUTPUT" | grep -q "initSchema" && echo "$LIST_OUTPUT" | grep -q "initSchema2"; then
    echo "✅ Migration entries found in list."
  else
    echo "❌ Migration entries not found in list."
    exit 1
  fi

  if echo "$LIST_OUTPUT" | grep -q "|             |"; then
    echo "✅ Executed At is EMPTY before running migrations (expected)."
  else
    echo "❌ Executed At is NOT empty before running migrations!"
    exit 1
  fi
}

function run_metana_up() {
  echo "⬆️  Running metana up..."
  (cd "$TEMP_DIR" && metana up)
  echo "✅ Up migrations executed."
}

function run_metana_list_and_validate_executed_at() {
  echo "📋 Running metana list (after running migrations)..."

  LIST_OUTPUT=$(cd "$TEMP_DIR" && metana list)

  echo "$LIST_OUTPUT"

  if echo "$LIST_OUTPUT" | grep -q "[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]"; then
    echo "✅ Executed At timestamp is present after running migrations."
  else
    echo "❌ Executed At timestamp NOT present after running migrations!"
    exit 1
  fi
}

function main() {
  check_metana_binary
  create_temp_dir
  run_metana_init
  create_migrations
  run_metana_list_and_validate_empty_executed
  run_metana_up
  run_metana_list_and_validate_executed_at

  echo "🎉 LIST migration tests passed successfully!"
}

main "$@"
