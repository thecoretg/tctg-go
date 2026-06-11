THIS_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SCRIPTS  := $(THIS_DIR)../scripts

-include .env
export

REGION      ?= us-west-2
ROLE_NAME   ?= salesforce-query-role
LAMBDA_NAME ?= salesforce-query
ARCH        := arm64
RUNTIME     := provided.al2023
TIMEOUT     ?= 300
MEMORY_SIZE ?= 512
LAMBDA_ENV_VARS = Variables={SALESFORCE_CLIENT_ID=$(SALESFORCE_CLIENT_ID),SALESFORCE_CLIENT_SECRET=$(SALESFORCE_CLIENT_SECRET),SALESFORCE_COMPANY_URL_NAME=$(SALESFORCE_COMPANY_URL_NAME),API_SECRET=$(API_SECRET)}

export REGION ROLE_NAME LAMBDA_NAME ARCH RUNTIME TIMEOUT MEMORY_SIZE LAMBDA_ENV_VARS

.PHONY: all iam build lambda url run clean

all: iam lambda url

iam:
	@$(SCRIPTS)/iam.sh

build:
	@$(SCRIPTS)/build.sh

lambda: build
	@$(SCRIPTS)/lambda.sh

url:
	@$(SCRIPTS)/url.sh

run:
	@go run ./cmd/local $(ARGS)

clean:
	@rm -rf dist
