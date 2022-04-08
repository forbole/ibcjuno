VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

all: ci-lint ci-test install

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

###############################################################################
###                                  Build                                  ###
###############################################################################

LD_FLAGS = -X github.com/forbole/ibcjuno/cmd.Version=$(VERSION) \
 	-X github.com/forbole/ibcjuno/cmd.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building ibcjuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/ibcjuno.exe  ./cmd/ibcjuno
else
	@echo "building ibcjuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/ibcjuno ./cmd/ibcjuno
endif

install: go.sum
	@echo "installing ibcjuno binary..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/ibcjuno

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	golangci-lint run --out-format=tab --timeout=10m

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0 --timeout=10m
.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs goimports -w -local github.com/forbole/ibcjuno
.PHONY: format

###############################################################################
###                             Tests & CI                                  ###
###############################################################################

coverage:
	@echo "viewing test coverage..."
	@go tool cover --html=coverage.out

ci-test:
	@echo "executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.out ./... 

ci-lint:
	@echo "running GolangCI-Lint..."
	@GO111MODULE=on golangci-lint run
	@echo "formatting..."
	@find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -d -s
	@echo "verifying modules..."
	@go mod verify

clean:
	rm -f tools-stamp ./build/**

.PHONY: install build ci-test ci-lint coverage clean
