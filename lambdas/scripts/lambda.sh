#!/usr/bin/env bash
set -euo pipefail

REGION="${REGION:-us-west-2}"
ROLE_NAME="${ROLE_NAME:-}"
LAMBDA_NAME="${LAMBDA_NAME:-}"
ARCH="${ARCH:-arm64}"
RUNTIME="${RUNTIME:-provided.al2023}"
TIMEOUT="${TIMEOUT:-300}"
MEMORY_SIZE="${MEMORY_SIZE:-512}"

if [[ -z "$ROLE_NAME" ]]; then
  echo "Error: ROLE_NAME is not set" >&2
  exit 1
fi

if [[ -z "$LAMBDA_NAME" ]]; then
  echo "Error: LAMBDA_NAME is not set" >&2
  exit 1
fi

if [[ -z "${LAMBDA_ENV_VARS:-}" ]]; then
  echo "Error: LAMBDA_ENV_VARS is not set" >&2
  exit 1
fi

ROLE_ARN=$(aws iam get-role --role-name "$ROLE_NAME" --query Role.Arn --output text)

if aws lambda get-function --function-name "$LAMBDA_NAME" --region "$REGION" &>/dev/null; then
  echo "Updating $LAMBDA_NAME..."

  aws lambda update-function-code \
    --function-name "$LAMBDA_NAME" \
    --zip-file "fileb://dist/bootstrap.zip" \
    --architectures "$ARCH" \
    --region "$REGION" \
    --output json > /dev/null

  aws lambda wait function-updated --function-name "$LAMBDA_NAME" --region "$REGION"

  aws lambda update-function-configuration \
    --function-name "$LAMBDA_NAME" \
    --environment "$LAMBDA_ENV_VARS" \
    --timeout "$TIMEOUT" \
    --memory-size "$MEMORY_SIZE" \
    --region "$REGION" \
    --output json > /dev/null
else
  echo "Creating $LAMBDA_NAME..."

  aws lambda create-function \
    --function-name "$LAMBDA_NAME" \
    --runtime "$RUNTIME" \
    --role "$ROLE_ARN" \
    --handler bootstrap \
    --zip-file "fileb://dist/bootstrap.zip" \
    --architectures "$ARCH" \
    --timeout "$TIMEOUT" \
    --memory-size "$MEMORY_SIZE" \
    --environment "$LAMBDA_ENV_VARS" \
    --region "$REGION" \
    --output json > /dev/null

  aws lambda wait function-active --function-name "$LAMBDA_NAME" --region "$REGION"
fi

echo "Done: $LAMBDA_NAME"
