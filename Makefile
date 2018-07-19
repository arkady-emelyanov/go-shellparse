.PHONY: format
format:
	gofmt -d -w -s -e .

.PHONY: coverage
coverage:
	go test -coverprofile=cover.out; go tool cover -html=cover.out
