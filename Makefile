# Directory to put `go install`ed binaries in.
export GOBIN ?= $(shell pwd)/bin

GO_FILES := $(shell \
	find . '(' -path '*/.*' -o -path './vendor' ')' -prune \
	-o -name '*.go' -print | cut -b3-)

.PHONY: bench
bench:
	go test -bench=. ./...

bin/golint: tools/go.mod
	@cd tools && go install golang.org/x/lint/golint

bin/staticcheck: tools/go.mod
	@cd tools && go install honnef.co/go/tools/cmd/staticcheck

.PHONY: build
build:
	go build ./...


.PHONY: run
run:
	go fmt ./...
	go build ./...
	go run main.go -c benthos.yaml

.PHONY: cover
cover:
	go test -coverprofile=cover.out -coverpkg=./... -v ./...
	go tool cover -html=cover.out -o cover.html
	go test -v ./... -covermode=count -coverprofile=coverage.out
	go tool cover -func=coverage.out -o=coverage.out

.PHONY: gofmt
gofmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -e -s -l $(GO_FILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || (echo "gofmt failed:" | cat - $(FMT_LOG) && false)

.PHONY: golint
golint: bin/golint
	@$(GOBIN)/golint -set_exit_status ./...

.PHONY: lint
lint: gofmt golint staticcheck

.PHONY: staticcheck
staticcheck: bin/staticcheck
	@$(GOBIN)/staticcheck ./...

.PHONY: test
test:
	go test -race ./...