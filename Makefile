# Timeout Utility Makefile
# All builds are statically linked for maximum portability

# Variables
BINARY_NAME := timeout
GO_FILES := timeout.go
TEST_FILES := timeout_test.go integration_test.go
VERSION := 1.0
BUILD_DIR := build
INSTALL_DIR := /usr/local/bin

# Build flags for static linking
STATIC_FLAGS := -ldflags "-extldflags '-static'"
STATIC_RELEASE_FLAGS := -ldflags "-X main.version=$(VERSION) -s -w -extldflags '-static'"

# Default target
.PHONY: all
all: build

# Help target
.PHONY: help
help:
	@echo "Timeout Utility Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  all     - Build the binary (default)"
	@echo "  build   - Build the binary (statically linked)"
	@echo "  clean   - Clean build artifacts"
	@echo "  status  - Show build status"
	@echo "  help    - Show this help"

# Build the binary (statically linked)
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) (static)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(STATIC_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(GO_FILES)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f timeout_test

# Check build status
.PHONY: status
status:
	@echo "Timeout Build Status:"
	@echo "===================="
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "✓ Binary: $(BUILD_DIR)/$(BINARY_NAME)"; \
		ls -lh $(BUILD_DIR)/$(BINARY_NAME); \
		file $(BUILD_DIR)/$(BINARY_NAME); \
	else \
		echo "✗ Binary: $(BUILD_DIR)/$(BINARY_NAME) (missing)"; \
	fi
