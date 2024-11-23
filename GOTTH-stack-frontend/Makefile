.PHONY: all build run test lint lint-fix clean

# Build the application
build:
	go build -o ./dist/server ./cmd/server

# Run the application using Air
run:
	air

# Run tests
test:
	go test ./...

# Run Revive linter
lint:
	revive -config revive.toml -formatter friendly ./...

# Run Revive linter with auto-fix (where possible)
lint-fix:
	revive -config revive.toml -formatter stylish -fix ./...

# Install dependencies
deps:
	go mod tidy
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/mgechev/revive@latest
	go install github.com/cosmtrek/air@latest
	npm install

# Clean build artifacts
clean:
	rm -rf tmp/
	go clean

# Default target
all: deps build
