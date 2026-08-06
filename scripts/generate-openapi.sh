#!/bin/bash
set -e

echo "Generating code from OpenAPI specification..."

# Pinned to match the Dockerfile: v2.8.0+ requires Go >= 1.25
OAPI_CODEGEN_VERSION=v2.7.0

# Check if oapi-codegen is installed
if ! command -v oapi-codegen >/dev/null 2>&1; then
    echo "Installing oapi-codegen ${OAPI_CODEGEN_VERSION}..."
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}
fi

# Create output directory
mkdir -p internal/generated

# Generate code
oapi-codegen -config api/codegen-config.yaml api/openapi.yaml

echo "✓ Code generated in internal/generated/api.gen.go"
