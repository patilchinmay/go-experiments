# OAPI Codegen and GOTTH Stack

A modern web application demonstrating spec-first development using the **GOTTH** stack:

- **G**o (backend language)
- **O**penAPI (API specification)
- **T**empl (HTML templating)
- **T**ailwindCSS (styling)
- **H**TMX (frontend interactivity)

## Table of Contents

- [OAPI Codegen and GOTTH Stack](#oapi-codegen-and-gotth-stack)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Architecture](#architecture)
  - [Project Structure](#project-structure)
  - [Key Features](#key-features)
    - [🔄 **Dual Response Format**](#-dual-response-format)
    - [📋 **OpenAPI-First Development**](#-openapi-first-development)
    - [⚡ **Modern Frontend Experience**](#-modern-frontend-experience)
    - [🏗️ **Clean Architecture**](#️-clean-architecture)
  - [Technology Stack](#technology-stack)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Running the Application](#running-the-application)
  - [Development Workflow](#development-workflow)
    - [Code Generation](#code-generation)
    - [Linting and Formatting](#linting-and-formatting)
    - [Live Development](#live-development)
  - [API Usage](#api-usage)
    - [Web Interface](#web-interface)
    - [REST API](#rest-api)
  - [Code Generation](#code-generation-1)
    - [1. OpenAPI Code Generation](#1-openapi-code-generation)
    - [2. Templ Template Generation](#2-templ-template-generation)
    - [3. TailwindCSS Generation](#3-tailwindcss-generation)
  - [Project Structure Details](#project-structure-details)
    - [Handler Logic (`internal/handlers/petstore.go`)](#handler-logic-internalhandlerspetstorego)
    - [View Templates (`internal/views/*.templ`)](#view-templates-internalviewstempl)
    - [Static Assets (`public/`)](#static-assets-public)
    - [Configuration Files](#configuration-files)
  - [References](#references)

## Overview

This project demonstrates **spec-first development** where the OpenAPI specification serves as the single source of truth for both API contracts and code generation. The application showcases how a single Go backend can serve both traditional REST API clients and modern web interfaces using HTMX.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Browser   │    │   CLI Client    │    │   API Client    │
│     (HTMX)      │    │  (Generated)    │    │   (Postman)     │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ HTML Responses       │ JSON Responses       │ JSON Responses
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      Go Server            │
                    │  ┌─────────────────────┐  │
                    │  │   Echo Router       │  │
                    │  │  ┌───────────────┐  │  │
                    │  │  │   Handlers    │  │  │
                    │  │  │ (Check HTMX   │  │  │
                    │  │  │  Headers)     │  │  │
                    │  │  └───────────────┘  │  │
                    │  └─────────────────────┘  │
                    └───────────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │     OpenAPI Spec          │
                    │  (petstore-expanded.yaml) │
                    └───────────────────────────┘
```

## Project Structure

```
oapi-codegen-gotth/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── handlers/
│   │   └── petstore.go            # Business logic handlers
│   ├── renderers/
│   │   └── home.go                # Page renderers
│   └── views/
│       ├── *.templ                # Templ templates
│       └── *_templ.go             # Generated Go code from templates
├── pkg/
│   └── spec/
│       ├── codegen-config/        # oapi-codegen configuration
│       └── generated/             # Generated Go code from OpenAPI
├── public/
│   ├── assets.go                  # Embedded static assets
│   └── css/
│       ├── input.css              # Tailwind input
│       └── output.css             # Generated Tailwind CSS
├── Makefile                       # Build and development commands
├── go.mod                         # Go module definition
├── package.json                   # Node.js dependencies (Tailwind)
└── tailwind.config.js             # Tailwind configuration
```

## Key Features

### 🔄 **Dual Response Format**

- **HTMX Requests**: Return HTML fragments for seamless page updates
- **API Requests**: Return JSON for traditional REST clients
- **Single Codebase**: Same handlers serve both formats based on request headers

### 📋 **OpenAPI-First Development**

- OpenAPI specification defines the API contract
- Server and client code generated automatically
- Request/response validation built-in
- Documentation stays in sync with implementation

### ⚡ **Modern Frontend Experience**

- **HTMX**: Dynamic interactions without JavaScript frameworks
- **Templ**: Type-safe HTML templates in Go
- **TailwindCSS**: Utility-first CSS framework
- **Real-time Updates**: Live pet list updates and modal interactions

### 🏗️ **Clean Architecture**

- Separation of concerns with clear layers
- Embedded static assets for single binary deployment
- Structured logging with contextual information
- Middleware for validation, logging, and recovery

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Backend** | Go 1.23+ | Server-side logic and API |
| **Web Framework** | Echo v4 | HTTP routing and middleware |
| **API Spec** | OpenAPI 3.0 | API contract definition |
| **Code Generation** | oapi-codegen | Generate Go code from OpenAPI |
| **Templates** | Templ | Type-safe HTML templating |
| **Frontend** | HTMX | Dynamic web interactions |
| **Styling** | TailwindCSS | Utility-first CSS |
| **Static Assets** | Go embed | Embedded file system |
| **Logging** | slog + tint | Structured logging |

## Getting Started

### Prerequisites

- Go 1.23+
- Node.js 16+: for TailwindCSS
- Make: Build automation tool

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd oapi-codegen-gotth
   ```

2. **Install dependencies**

   ```bash
   make deps
   ```

   This command installs:
   - Go dependencies (`go mod tidy`)
   - oapi-codegen CLI tool
   - templ CLI tool
   - revive linter
   - air live reload tool
   - Node.js dependencies (TailwindCSS)

### Running the Application

1. **Generate code and start the server**

   ```bash
   make run
   ```

2. **Open your browser**
   - Web Interface: <http://localhost:3000/ui>
   - API Base: <http://localhost:3000>

## Development Workflow

### Code Generation

```bash
# Generate all code (OpenAPI + Templ + CSS)
make generate

# Individual generation steps:
# 1. Generate Go types from OpenAPI spec
# 2. Generate server handlers from OpenAPI spec
# 3. Generate client code from OpenAPI spec
# 4. Compile Tailwind CSS
# 5. Generate Go code from Templ templates
```

### Linting and Formatting

```bash
# Run linters and formatters
make lint
```

### Live Development

```bash
# Use air for live reload during development
air
```

## API Usage

### Web Interface

The web interface demonstrates modern HTMX patterns:

**Pet Management Features:**

- ✅ Add new pets with real-time form submission
- ✅ View pet list with dynamic loading
- ✅ Filter pets by tags with live search
- ✅ View pet details in modal dialogs
- ✅ Delete pets with confirmation prompts
- ✅ Real-time UI updates without page refreshes

**HTMX Interactions:**

```html
<!-- Auto-load pet list on page load -->
<div hx-get="/pets" hx-trigger="load" hx-target="#pet-list">

<!-- Live search with debouncing -->
<input hx-get="/pets" hx-trigger="keyup changed delay:500ms" hx-target="#pet-list">

<!-- Form submission with JSON encoding -->
<form hx-post="/pets" hx-ext="json-enc" hx-target="#pet-list" hx-swap="beforeend">

<!-- Delete with confirmation -->
<button hx-delete="/pets/123" hx-confirm="Are you sure?">
```

### REST API

Standard REST endpoints for API clients:

```bash
# Get all pets
curl http://localhost:3000/pets

# Get pets with filtering
curl "http://localhost:3000/pets?tags=dog&limit=10"

# Get specific pet
curl http://localhost:3000/pets/123

# Create new pet
curl -X POST http://localhost:3000/pets \
  -H "Content-Type: application/json" \
  -d '{"name": "Buddy", "tag": "dog"}'

# Delete pet
curl -X DELETE http://localhost:3000/pets/123
```

## Code Generation

The project uses multiple code generation tools:

### 1. OpenAPI Code Generation

```yaml
# pkg/spec/codegen-config/types.cfg.yaml
package: generated
output: pkg/spec/generated/types.go
generate:
  models: true

# pkg/spec/codegen-config/server.cfg.yaml
package: generated
output: pkg/spec/generated/server.go
generate:
  echo-server: true

# pkg/spec/codegen-config/client.cfg.yaml
package: generated
output: pkg/spec/generated/client.go
generate:
  client: true
```

### 2. Templ Template Generation

```go
// internal/views/home.templ -> internal/views/home_templ.go
templ HomePage() {
    @Layout() {
        <div class="space-y-6">
            <!-- Template content -->
        </div>
    }
}
```

### 3. TailwindCSS Generation

```css
/* public/css/input.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

/* Generates -> public/css/output.css */
```

## Project Structure Details

### Handler Logic (`internal/handlers/petstore.go`)

- Implements the generated `ServerInterface`
- Checks `HX-Request` header to determine response format
- Returns HTML fragments for HTMX or JSON for API clients
- Thread-safe in-memory pet storage with mutex

### View Templates (`internal/views/*.templ`)

- **Layout**: Base HTML structure with navigation
- **Home**: Main page with pet form and list
- **Pet**: Reusable pet card and detail components
- Type-safe with Go integration

### Static Assets (`public/`)

- Embedded using Go's `embed` package
- Single binary deployment with all assets included
- Served via Echo's static file handler

### Configuration Files

- **Makefile**: Development commands and build automation
- **revive.toml**: Go linting configuration
- **tailwind.config.js**: CSS framework configuration
- **.pre-commit-config.yaml**: Git hooks for code quality

## References

- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - OpenAPI Code Generator for Go
- [Templ](https://templ.guide/) - Type-safe HTML templating for Go
- [HTMX](https://htmx.org/) - High power tools for HTML
- [Echo](https://echo.labstack.com/) - High performance Go web framework
- [TailwindCSS](https://tailwindcss.com/) - Utility-first CSS framework
- [OpenAPI Specification](https://swagger.io/specification/) - API specification standard
- [Petstore Example](https://github.com/oapi-codegen/oapi-codegen/blob/main/examples/petstore-expanded/petstore-expanded.yaml) - Sample OpenAPI specification
