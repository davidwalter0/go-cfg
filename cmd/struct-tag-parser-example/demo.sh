#!/usr/bin/env bash

echo "=== Struct Tag Parser Integration Demo ==="
echo
echo "This demonstrates how go-cfg automatically generates help text"
echo "from usage tags that would be extracted from comments."
echo

echo "1. Show help output (generated from usage tags):"
echo "   go run . --help"
echo
go run . --help 2>&1 | head -20
echo "   ... (truncated for brevity)"
echo

echo "2. Run with basic configuration:"
echo "   go run . --service-name \"demo-api\" --debug"
echo
go run . --service-name "demo-api" --debug
echo

echo "3. Run with nested struct flags:"
echo "   go run . --database-host \"prod-db\" --server-port 9000 --monitoring-log-level \"warn\""
echo
go run . --database-host "prod-db" --server-port 9000 --monitoring-log-level "warn"
echo

echo "4. Environment variable integration:"
echo "   export SERVICE_NAME=\"env-api\" DATABASE_HOST=\"env-db\" && go run ."
echo
export SERVICE_NAME="env-api" 
export DATABASE_HOST="env-db"
go run .
echo

echo "=== Key Benefits ==="
echo "✅ Automatic help text from usage tags (extracted from comments)"
echo "✅ Nested struct support with kebab-case flags"
echo "✅ Environment variable integration with SCREAMING_SNAKE_CASE"
echo "✅ Default values from struct tags"
echo "✅ Type-safe flag parsing with go-cfg"
echo
echo "In a real implementation:"
echo "- Comments would be automatically converted to usage tags via go:generate"
echo "- No manual struct tag writing required"
echo "- Single source of truth: write documentation once in comments"