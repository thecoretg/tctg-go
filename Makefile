.PHONY: test test-td test-rewst test-sf test-iru test-entra

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

test-entra:
	go test -v ./entra/...

