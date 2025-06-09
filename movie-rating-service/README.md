# 🎬 Movie Rating Service

A clean architecture, domain-driven Go backend for rating movies, with JWT authentication, PostgreSQL, caching, and robust testability.

---

## 🚀 Features

* **User Registration & JWT Login**
* **Graceful Shutdown** (handles SIGINT/SIGTERM, finishes active requests)
* **Movie CRUD** and search endpoints
* **Movie Rating** (user-to-movie, one rating per user/movie)
* **User Profile & Rated Movies**
* **Caching Decorators** (in-memory, TTL-based for movies and ratings)
* **PostgreSQL** (GORM, auto-migration)
* **Swagger/OpenAPI** documentation (`/swagger/index.html`)
* **Prometheus Metrics** (`/metrics`)
* **Domain-Driven Structure** (DDD, Clean Architecture)
* **Configurable via YAML or ENV**
* **Production-Ready Lint** (`golangci-lint`)

---

## 🗂️ Project Structure

```
movie-rating-service/
│
├── config/                 # Config parsing (YAML, ENV)
├── docs/                   # Swagger docs (auto-generated)
├── internal/
│   ├── application/        # Use-case (controller,service) logic
│   ├── common/             # Shared helpers, utilities, error handling, etc.
│   ├── domain/             # Pure domain entities & interfaces
│   └── infrastructure/     # DB, repo implementations, cache
├── mocks/                  # Generated mocks (e.g., for testing)
├── .golangci.yml           # Linter configuration
├── docker-compose.yml      # Docker Compose config for local development
├── Dockerfile              # App Docker image build instructions
├── go.mod                  # Go module definition
├── main.go                 # Application entry point (if not under /cmd)
└── README.md               # Project documentation
```

---

## ⚡️ Quick Start (Development)

### 1. Clone the repo

```sh
git clone https://github.com/EmreAyberk/thermondo.git
cd movie-rating-service
```

### 2. Set up Database and App (via Docker Compose)

```sh
docker-compose up -d
```

---

## 📚 API Documentation

* **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
* **OpenAPI JSON:** [http://localhost:8080/swagger/doc.json](http://localhost:8080/swagger/doc.json)

---

## 📈 Monitoring

* **Prometheus Metrics:** [http://localhost:8080/metrics](http://localhost:8080/metrics)

---

## 🧩 Key Endpoints

| Method | Endpoint            | Description                  |
| ------ |---------------------|------------------------------|
| POST   | `/login`            | User JWT login               |
| POST   | `/user`             | Create user                  |
| GET    | `/users/:id`        | Get user profile (auth)      |
| GET    | `/movies`           | List all movies              |
| POST   | `/movies`           | Add a new movie (admin/auth) |
| GET    | `/movies/:id`       | Movie details                |
| POST   | `/ratings`          | Rate a movie (auth)          |
| GET    | `/ratings/user/:id` | Get ratings by user          |

> For the full API and request/response schema, check `/swagger/index.html`.

---

## 🧪 Testing

* Unit & integration tests are in the `test/` folder.
* Run tests:

  ```sh
  go test ./...
  ```

---

## 🧹 Linting & Code Quality

* Use [golangci-lint](https://github.com/golangci/golangci-lint) for style and static analysis.
* Run:

  ```sh
  golangci-lint run
  ```

---

## 🎨 Architecture

* DDD layered (domain/application/infrastructure/interfaces)
* Service/repo pattern with interfaces for testability
* Decorator pattern for caching
* Middleware for JWT authentication
* Configurable, minimal, and robust

For details, see the included Excalidraw diagram (https://excalidraw.com/#json=-HZmmTZG3ueNTfj5OLshh,hZ7jVlM6zlX3jdKwtzsMWA).

---
