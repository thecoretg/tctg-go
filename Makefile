.PHONY: test test-td test-rewst

test:
	go test -v ./...

test-td:
	go test -v ./threatdown/...

test-rewst:
	go test -v ./rewst/...
