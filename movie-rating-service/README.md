# 🎬 Movie Rating Service

A clean architecture, domain-driven Go backend for rating movies, with JWT authentication, PostgreSQL, caching, and
robust testability.

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
├── .github/                        # GitHub configuration (actions)
├── config/                         # Config parsing (ENV)
├── docs/                           # Swagger docs (auto-generated)
├── internal/
│   ├── application/                # Use-case logic
│   │   ├── controller/             # HTTP handlers/controllers
│   │   ├── middleware/             # Auth/JWT and other middleware
│   │   ├── models/                 # Request/response DTOs
│   │   ├── service/                # Application services
│   │   └── validator/              # Request validators
│   ├── common/
│   │   └── errorhandler.go         # Error handling utilities
│   ├── domain/                     # Domain entities
│   └── infrastructure/
│       ├── db/
│       │   ├── seeder/             # DB seeder
│       │   └── postgres.go         # DB connection/config
│       └── repository/             # Repository implementations (GORM)
├── k8s-manifest/                   # Kubernetes deployment and service YAMLs
├── mocks/                          # Generated mocks (for testing)
├── .golangci.yml                   # Linter configuration
├── docker-compose.yml              # Docker Compose config for local development
├── Dockerfile                      # App Docker image build instructions
├── go.mod                          # Go module definition
├── main.go                         # Application entry point
├── Makefile                        # Makefile for commands (if present)
└── README.md                       # Project documentation
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

## 🌱 Database Seeder

To populate your database with demo users, movies, and ratings, you can manually run the seeder at any time:

```sh
docker-compose run --rm movie-rating-service-app go run main.go seeder
```

- This will execute the seeding logic and then exit.
- **Seeder does not run automatically** on `docker-compose up`—you must invoke it explicitly when you want to re-seed
  your data.

## 📚 API Documentation

* **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
* **OpenAPI JSON:** [http://localhost:8080/swagger/doc.json](http://localhost:8080/swagger/doc.json)

---

## 📈 Monitoring

* **Prometheus Metrics:** [http://localhost:8080/metrics](http://localhost:8080/metrics)

---

## 🧩 Key Endpoints
### Users

| Method | Endpoint       | Description              |
|--------|---------------|-------------------------|
| POST   | `/login`      | User JWT login          |
| POST   | `/user`       | Create user             |
| GET    | `/users/:id`  | Get user profile (auth) |

---

### Movies

| Method | Endpoint        | Description                   |
|--------|----------------|------------------------------|
| GET    | `/movies`      | List all movies              |
| POST   | `/movies`      | Add a new movie (admin/auth) |
| GET    | `/movies/:id`  | Movie details                |

---

### Ratings

| Method | Endpoint        | Description                                 |
|--------|----------------|---------------------------------------------|
| POST   | `/rating`      | Create a new rating (auth required)         |
| PATCH  | `/rating`      | Update a rating (auth required)             |
| DELETE | `/rating`      | Delete a rating (auth required)             |
| GET    | `/rating/user` | List all ratings by the authenticated user  |

---

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

For details, see the included Excalidraw
diagram (https://excalidraw.com/#json=-HZmmTZG3ueNTfj5OLshh,hZ7jVlM6zlX3jdKwtzsMWA).

---
