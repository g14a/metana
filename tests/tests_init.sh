#!/bin/bash

set -e

echo "🔵 Starting INIT migration tests..."

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

function validate_default_migrations_structure() {
  echo "🔍 Validating default migrations structure..."

  if [ ! -d "$TEMP_DIR/migrations/scripts" ]; then
    echo "❌ migrations/scripts directory not created!"
    exit 1
  fi

  echo "✅ Default migrations/scripts present."
}

function main() {
  check_metana_binary
  create_temp_dir

  echo "==============================="
  echo "▶️  Testing default init..."
  run_metana_init_default
  validate_default_migrations_structure

  echo "==============================="
  echo "🎉 ALL INIT tests passed successfully!"
}

main "$@"
