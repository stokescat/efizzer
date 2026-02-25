
# Applications
APPNAME_MANAGER := efizzer-manager
APPNAME_ORACLE  := efizzer-oracle


BINDIR:= bin
CMDDIR:= ./cmd


GO      := go
GOFLAGS := -mod=readonly


PACKAGES := ./internal/rawcov \
            ./internal/efi

.PHONY: all build test clean help

all: test build

build: $(BINDIR)/$(APPNAME_MANAGER) $(BINDIR)/$(APPNAME_ORACLE)

$(BINDIR)/$(APPNAME_MANAGER):
	@mkdir -p $(BINDIR)
	@echo "Building $(APPNAME_MANAGER)..."
	$(GO) build $(GOFLAGS) -o $(BINDIR)/$(APPNAME_MANAGER) $(CMDDIR)/$(APPNAME_MANAGER)

$(BINDIR)/$(APPNAME_ORACLE):
	@mkdir -p $(BINDIR)
	@echo "Building $(APPNAME_ORACLE)..."
	$(GO) build $(GOFLAGS) -o $(BINDIR)/$(APPNAME_ORACLE) $(CMDDIR)/$(APPNAME_ORACLE)


test:
	@echo "Running tests..."
	$(GO) test -v $(PACKAGES)

clean:
	@echo "Cleaning up..."
	rm -rf $(BINDIR)

help:
	@echo "Usage:"
	@echo "  make build   - Build all binaries"
	@echo "  make test    - Run all tests"
	@echo "  make clean   - Remove binaries"
	@echo "  make all     - Run tests and build"
