# OAPI Codegen and GOTTH stack

- [OAPI Codegen and GOTTH stack](#oapi-codegen-and-gotth-stack)
- [Description](#description)

# Description

This repo demonstrates spec first development.

Key points about this setup:

- The `OpenAPI` spec remains single source of truth.
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen/tree/main) generates both server and client code.
- Handlers check for `HTMX` headers to determine response type.
  - Requests from CLI client will receive JSON responses.
  - Requests from Web clients will receive HTML responses.
- `templ` components are used for HTML rendering.
- The CLI uses the generated client code.
- The web interface uses `HTMX` for dynamic updates.
- Single backend serves both Web and CLI clients.

For this example we will use the sample spec [petstore-expanded](https://github.com/oapi-codegen/oapi-codegen/blob/main/examples/petstore-expanded/petstore-expanded.yaml) from oapi.
