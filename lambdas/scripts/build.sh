#!/usr/bin/env bash
set -euo pipefail

ARCH="${ARCH:-arm64}"

echo "Building..."
mkdir -p dist
GOOS=linux GOARCH=$ARCH go build -o dist/bootstrap ./cmd/lambda
(cd dist && zip -q bootstrap.zip bootstrap)
echo "Built dist/bootstrap.zip"
