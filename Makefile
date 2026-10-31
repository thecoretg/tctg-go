.PHONY: test test-td test-sf test-entra

test:
	go test -v ./...

test-td:
	go test -v ./threatdown/...

test-sf:
	go test -v ./salesforce/...

test-entra:
	go test -v ./entra/...

