.PHONY: test test-td test-rewst test-sf test-iru

test:
	go test -v ./...

test-td:
	go test -v ./threatdown/...

test-rewst:
	go test -v ./rewst/...

test-sf:
	go test -v ./salesforce/...

test-iru:
	go test -v ./iru/...

# --- Lambda targets (delegate to each lambda's Makefile) ---
# Usage: make lambda-<name> [ARGS="<target>"]
# Example: make lambda-threatdown-site-list ARGS=lambda

LAMBDA_DIRS := $(wildcard lambdas/*/Makefile)
LAMBDAS     := $(patsubst lambdas/%/Makefile,%,$(LAMBDA_DIRS))

.PHONY: lambdas $(addprefix lambda-,$(LAMBDAS))

lambdas: $(addprefix lambda-,$(LAMBDAS))

$(addprefix lambda-,$(LAMBDAS)): lambda-%:
	$(MAKE) -C lambdas/$* $(ARGS)
