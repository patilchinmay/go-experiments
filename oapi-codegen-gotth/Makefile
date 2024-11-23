# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/mgechev/revive@latest
	go install github.com/air-verse/air@latest
	npm install

# Code Generation
.PHONY: generate
generate:
	oapi-codegen -config pkg/spec/codegen-config/types.cfg.yaml pkg/spec/petstore-expanded.yaml
	oapi-codegen -config pkg/spec/codegen-config/server.cfg.yaml pkg/spec/petstore-expanded.yaml
	oapi-codegen -config pkg/spec/codegen-config/client.cfg.yaml pkg/spec/petstore-expanded.yaml

# Run Revive linter
.PHONY: lint
lint:
	revive -config revive.toml -formatter stylish ./...
