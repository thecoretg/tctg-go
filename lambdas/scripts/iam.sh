#!/usr/bin/env bash
set -euo pipefail

ROLE_NAME="${ROLE_NAME:-}"

if [[ -z "$ROLE_NAME" ]]; then
  echo "Error: ROLE_NAME is not set" >&2
  exit 1
fi

if aws iam get-role --role-name "$ROLE_NAME" &>/dev/null; then
  echo "IAM role already exists: $ROLE_NAME"
  exit 0
fi

echo "Creating IAM role: $ROLE_NAME"
aws iam create-role \
  --role-name "$ROLE_NAME" \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [{
      "Effect": "Allow",
      "Principal": { "Service": "lambda.amazonaws.com" },
      "Action": "sts:AssumeRole"
    }]
  }' \
  --output json > /dev/null

aws iam attach-role-policy \
  --role-name "$ROLE_NAME" \
  --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

echo "Waiting for role to propagate..."
sleep 10

echo "Created: $(aws iam get-role --role-name "$ROLE_NAME" --query Role.Arn --output text)"
