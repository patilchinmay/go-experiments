##############################################################
### Code Generation
##############################################################

.PHONY: generate-types
generate-types:
	@go tool oapi-codegen -config ./api/codegen-config/types.config.yaml ./api/spec/spec.yaml

.PHONY: generate-server
generate-server:
	@go tool oapi-codegen -config ./api/codegen-config/server.config.yaml ./api/spec/spec.yaml

.PHONY: generate
generate: generate-types generate-server

##############################################################
### Run
##############################################################

.PHONY: run
run:
	@go run ./cmd/api

