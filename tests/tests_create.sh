#!/bin/bash

set -e

echo "🔵 Starting INIT + CREATE migration tests..."

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
  cd $TEMP_DIR
}

function run_metana_init_default() {
  echo "🚀 Running metana init..."
  metana init
}

function validate_migrations_structure() {
  echo "🔍 Validating default migrations structure..."

  if [ ! -d "$TEMP_DIR/migrations/scripts" ]; then
    echo "❌ migrations/scripts directory not created!"
    exit 1
  fi

  echo "✅ migrations/scripts present."
}

function run_metana_create_default_dir() {
  echo "🚀 Running metana create initSchema (default dir)..."

  (cd "$TEMP_DIR" && metana create initSchema)

  CREATED_FILE=$(find "$TEMP_DIR/migrations/scripts" -type f -name "*_initSchema.go" || true)

  if [ -z "$CREATED_FILE" ]; then
    echo "❌ Migration file for initSchema not created in default migrations dir!"
    exit 1
  fi

  echo "✅ Migration created successfully in default migrations dir: $CREATED_FILE"
}

function main() {
  check_metana_binary
  create_temp_dir

  echo "==============================="
  echo "▶️  Testing metana init..."
  run_metana_init_default
  validate_migrations_structure

  echo "==============================="
  echo "▶️  Testing metana create in default dir..."
  run_metana_create_default_dir

  echo "==============================="
  echo "🎉 INIT + CREATE tests passed successfully!"
}

main "$@"
