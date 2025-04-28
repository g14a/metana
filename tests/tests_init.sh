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
}

function run_metana_init_default() {
  echo "🚀 Running metana init (default dir)..."
  metana init --dir "$TEMP_DIR/migrations"
}

function validate_default_migrations_structure() {
  echo "🔍 Validating default migrations structure..."

  if [ ! -d "$TEMP_DIR/migrations/scripts" ]; then
    echo "❌ migrations/scripts directory not created!"
    exit 1
  fi

  if [ ! -f "$TEMP_DIR/.metana.yml" ]; then
    echo "❌ .metana.yml not created!"
    exit 1
  fi

  echo "✅ Default migrations/scripts and .metana.yml are present."
}

function run_metana_init_custom_dir() {
  CUSTOM_DIR="$TEMP_DIR/schema-migrations"
  echo "🚀 Running metana init with --dir=$CUSTOM_DIR..."
  metana init --dir "$CUSTOM_DIR"
}

function validate_custom_migrations_structure() {
  echo "🔍 Validating custom migrations structure..."

  if [ ! -d "$CUSTOM_DIR/scripts" ]; then
    echo "❌ Custom migrations/scripts directory not created!"
    exit 1
  fi

  if [ ! -f "$TEMP_DIR/.metana.yml" ]; then
    echo "❌ .metana.yml not created for custom dir!"
    exit 1
  fi

  echo "✅ Custom migrations/scripts and .metana.yml are present."
}

function main() {
  check_metana_binary
  create_temp_dir

  echo "==============================="
  echo "▶️  Testing default init..."
  run_metana_init_default
  validate_default_migrations_structure

  echo "==============================="
  echo "▶️  Testing init with custom --dir..."
  run_metana_init_custom_dir
  validate_custom_migrations_structure

  echo "==============================="
  echo "🎉 ALL INIT tests passed successfully!"
}

main "$@"
