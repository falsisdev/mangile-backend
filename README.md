# Mangile Backend

Low-latency, high-throughput backend engine for the [Mangile](https://github.com/falsisdev/mangile) platform.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)

## Overview

`mangile-backend` is a scalable content serving and indexing service written in Go. Built with standard library efficiency (`net/http`), clean architecture, and strict domain separation, it powers the backend capabilities of the Mangile ecosystem.

### Key Features

* **High-Performance Pipeline:** Leverages Go's native concurrent runtime and optimized `net/http` stack for minimal overhead.
* **Strict Domain Isolation:** Explicit segregation between distinct comic/media entities (e.g., `Manga` vs. `LightNovel` domain models).
* **Modern Architecture:** Standardized, layered modular structure adhering to Go best practices (`cmd/`, `internal/`).

---

## Architecture & Structure

```text
.
├── cmd/             # Application entry points (e.g., cmd/api/main.go)
└── internal/        # Private application and domain logic
    ├── handlers/    # Transport protocol handlers (HTTP/REST)
    ├── models/      # Core data models and schema definitions
    └── services/    # Business logic implementations & external integrations
```

---

## Tech Stack

| Category | Technology |
| :--- | :--- |
| **Language** | Go (Golang 1.22+) |
| **Communication** | REST (HTTP/1.1 & HTTP/2 via `net/http`) |
| **CMS & Data Sources** | Sanity CMS, AniList API, MyAnimeList (MAL) API |

---

## Getting Started

### Prerequisites

* [Go](https://go.dev/doc/install) (v1.22 or higher)
* Active [Sanity](https://sanity.io) project credentials

### Local Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/falsisdev/mangile-backend.git
   cd mangile-backend
   ```

2. **Configure environment variables:**
   
   Copy the example environment file and update it with your Sanity credentials:
   ```bash
   cp .env.example ./cmd/api/.env
   ```

3. **Run the server:**
   ```bash
   go run ./cmd/api
   ```

---

## Development

Execute tests across all packages:

```bash
go test -v -race ./...
```

---

## License

Distributed under the MIT License. See `LICENSE` for more information.