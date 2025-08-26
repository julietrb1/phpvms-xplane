.PHONY: build run test clean

# Build variables
BINARY_NAME=pxp
BUILD_DIR=build

# Go commands
GO=go
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOCLEAN=$(GO) clean

# Build the application
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/pxp

# Run the application
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run with specific environment variables for development
dev: build
	PHPVMS_BASE_URL=http://localhost:8080 \
	PHPVMS_API_KEY=your_api_key \
	UDP_BIND_HOST=0.0.0.0 \
	UDP_BIND_PORT=47777 \
	LOG_LEVEL=debug \
	$(BUILD_DIR)/$(BINARY_NAME)
