# Configuration variables
APP ?= o7k-opentelekomcloud-plugin

TAG ?= 0.1.0
ITERATION ?= 1
TARGET_ENV ?= dev

VERSION ?= $(TAG)-$(TARGET_ENV).$(ITERATION)

# Build metadata
BIN_VERSION ?= $(VERSION)
BIN_GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BIN_BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = -s -w \
	-X github.com/akyriako/o7k-opentelekomcloud-plugin/internal/version.Version=$(BIN_VERSION) \
	-X github.com/akyriako/o7k-opentelekomcloud-plugin/internal/version.Commit=$(BIN_GIT_COMMIT) \
	-X github.com/akyriako/o7k-opentelekomcloud-plugin/internal/version.BuildDate=$(BIN_BUILD_DATE)

build:
	@echo "Building $(APP) $(BIN_VERSION)..."
	@mkdir -p bin
	CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o bin/$(APP) \
		./main.go

PLUGIN_DIR ?= $(HOME)/.config/o7k/plugins

install: build
	@echo "Installing $(APP) to $(PLUGIN_DIR)..."
	@mkdir -p "$(PLUGIN_DIR)"
	@rm -f "$(PLUGIN_DIR)/$(APP)"
	@cp "bin/$(APP)" "$(PLUGIN_DIR)/$(APP)"
	@echo "Installed $(PLUGIN_DIR)/$(APP)"

run:
	go run ./main.go

test:
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin dist

release-tag:
	@test "$$(git branch --show-current)" = "main" || { echo "Error: release-tag must run on main"; exit 1; }
	@git diff --quiet && git diff --cached --quiet || { echo "Error: working tree has uncommitted changes"; exit 1; }
	@git fetch origin main
	@test "$$(git rev-list --count origin/main..HEAD)" -eq 0 || { echo "Error: main has unpushed commits"; exit 1; }
	@test "$$(git rev-list --count HEAD..origin/main)" -eq 0 || { echo "Error: main is behind origin/main"; exit 1; }
	git tag -a v$(TAG) -m "$(APP) v$(TAG)"
	git push origin v$(TAG)

.PHONY: build run test clean install release-tag