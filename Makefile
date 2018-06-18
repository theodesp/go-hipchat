.PHONY: format
format:
	@find . -type f -name "*.go*" -print0 | xargs -0 gofmt -s -w

.PHONY: debs
debs:
	GOPATH=$(GOPATH) go get -u ./...
	GOPATH=$(GOPATH) go get -u github.com/stretchr/testify

.PHONY: test
test:
	GOPATH=$(GOPATH) go test -race

.PHONY: bench
bench:
	GOPATH=$(GOPATH) go test -bench=. -check.b -benchmem

# Clean junk
.PHONY: clean
clean:
	GOPATH=$(GOPATH) go clean ./...