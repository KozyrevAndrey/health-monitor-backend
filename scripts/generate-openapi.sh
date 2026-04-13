#!/bin/bash
set -e

echo "Generating code from OpenAPI specification..."

# Check if oapi-codegen is installed
if ! command -v oapi-codegen >/dev/null 2>&1; then
    echo "Installing oapi-codegen..."
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
fi

# Create output directory
mkdir -p internal/generated

# Generate code
oapi-codegen -config api/codegen-config.yaml api/openapi.yaml

echo "✓ Code generated in internal/generated/api.gen.go"
