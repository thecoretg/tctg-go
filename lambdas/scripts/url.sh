#!/usr/bin/env bash
set -euo pipefail

REGION="${REGION:-us-west-2}"
LAMBDA_NAME="${LAMBDA_NAME:-}"

if [[ -z "$LAMBDA_NAME" ]]; then
  echo "Error: LAMBDA_NAME is not set" >&2
  exit 1
fi

if ! aws lambda get-function-url-config \
  --function-name "$LAMBDA_NAME" \
  --region "$REGION" &>/dev/null; then

  echo "Creating function URL for $LAMBDA_NAME..."
  aws lambda create-function-url-config \
    --function-name "$LAMBDA_NAME" \
    --auth-type NONE \
    --region "$REGION" \
    --output json > /dev/null
else
  echo "Function URL already exists for $LAMBDA_NAME"
fi

aws lambda add-permission \
  --function-name "$LAMBDA_NAME" \
  --statement-id "FunctionURLAllowPublicAccess" \
  --action lambda:InvokeFunctionUrl \
  --principal "*" \
  --function-url-auth-type NONE \
  --region "$REGION" \
  --output json > /dev/null 2>&1 || true

aws lambda add-permission \
  --function-name "$LAMBDA_NAME" \
  --statement-id "FunctionURLAllowInvokeFunction" \
  --action lambda:InvokeFunction \
  --principal "*" \
  --region "$REGION" \
  --output json > /dev/null 2>&1 || true

URL=$(aws lambda get-function-url-config \
  --function-name "$LAMBDA_NAME" \
  --region "$REGION" \
  --query FunctionUrl \
  --output text)

echo "URL: ${URL}"
