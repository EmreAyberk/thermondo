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

## 🗃️ Domain Models

### User

Represents an application user.

- `ID` *(uint, primary key)*: Unique user identifier.
- `Username` *(unique)*: Login name, must be unique.
- `Password`: Hashed password for authentication.
- `Name`, `Surname`, `Email`, `Phone`, `Address`: Profile information.
- `IsAdmin` *(bool)*: Set to `true` for admin users.

---

### Movie

Represents a movie available for rating.

- `ID` *(uint, primary key)*: Unique movie identifier.
- `Title`, `Description`, `Genre`, `Director`, `Year`: Movie metadata.
- `Rating` *(float64)*: Average rating (calculated).
- `RatingCount` *(int64)*: Number of ratings for this movie.

---

### Rating

Represents a rating given by a user to a movie.

- `ID` *(uint, primary key)*: Unique rating identifier.
- `UserID`: The user who gave the rating (foreign key).
- `MovieID`: The movie being rated (foreign key).
- `Score` *(float64)*: The rating score (e.g., 0–5).
- `Review`: Optional text review.
- **Composite Unique Index:**
    - There is a unique constraint on (`UserID`, `MovieID`) to ensure that **each user can only rate each movie once**.

---

**Relationships:**

- A `User` can rate many `Movies`.
- A `Movie` can be rated by many `Users`.
- The `Rating` table links users and movies, with each (user, movie) pair unique.

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

## 🌱 Data Seeder

This project includes a **data seeder** for quickly populating your database with demo data for local development and
testing.  
The seeder will insert sample users, movies, and ratings:

- **Users:**  
  3 demo users (`alice`, `bob`, `carol`), each with basic profile information. The password for all is the same, `1234`
  hashed version and only `carol` is admin.
- **Movies:**  
  3 movies from different genres and directors, with basic metadata.
- **Ratings:**  
  5 example ratings linking users and movies, each with a score and a review.

**How to run the seeder:**  
Use this command to execute the seeder (from your project root):

```sh
docker-compose run --rm movie-rating-service-app go run main.go seeder
```

- The seeder will log its progress and insert all sample data into your connected database.
- You can rerun the seeder any time to reset or re-populate demo data.

**Note:**  
The seeder does **not** run automatically with the app; you must invoke it manually when needed.

## 📚 API Documentation

* **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
* **OpenAPI JSON:** [http://localhost:8080/swagger/doc.json](http://localhost:8080/swagger/doc.json)

---

## 📈 Metric & Monitoring

* **Prometheus Metrics:** [http://localhost:8080/metrics](http://localhost:8080/metrics)
* **Fiber Monitoring UI:** [http://localhost:8080/monitor](http://localhost:8080/monitor)

---

## 🧩 Key Endpoints

### Users

| Method | Endpoint    | Description             |
|--------|-------------|-------------------------|
| POST   | `/login`    | User JWT login          |
| POST   | `/user`     | Create user             |
| GET    | `/user/:id` | Get user profile (auth) |

---

### Movies

| Method | Endpoint     | Description                  |
|--------|--------------|------------------------------|
| GET    | `/movie`     | List all movies              |
| POST   | `/movie`     | Add a new movie (admin/auth) |
| GET    | `/movie/:id` | Movie details                |

---

### Ratings

| Method | Endpoint            | Description                                |
|--------|---------------------|--------------------------------------------|
| POST   | `/movie/:id/rating` | Create a new rating (auth required)        |
| PATCH  | `/movie/:id/rating` | Update a rating (auth required)            |
| DELETE | `/movie/:id/rating` | Delete a rating (auth required)            |
| GET    | `/user/rating`      | List all ratings by the authenticated user |

---

## 🔐 Authentication & Authorization

All protected endpoints use JWT-based authentication, handled by custom middleware.

### User vs. Admin Middleware

- **UserHandler:**
    - Checks for a valid JWT in the `Authorization` header.
    - Validates the token and extracts user claims.
    - On success, attaches claims to the request context (accessible via `ctx.Locals("user")`).
    - Grants access to any authenticated user.

- **AdminHandler:**
    - Performs all checks of `UserHandler`.
    - Additionally verifies that the `isAdmin` claim is present and set to `true`.
    - Denies access (`401 Unauthorized`) if the user is not an admin.

**In summary:**

- Use `UserHandler` to protect routes accessible to any logged-in user.
- Use `AdminHandler` to restrict routes to admin users only.

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

## 🗄️ Caching

Caching is applied using the **decorator pattern** to wrap repository methods with in-memory, TTL-based cache layers.
This accelerates repeated read operations and reduces database load for frequently requested data.

- **Where is caching used?**
    - **Movies:** List and details are cached with per-item TTL.

- **How does it work?**
    - Decorator checks the in-memory cache before hitting the database.
    - Cache entries expire after a fixed TTL.
    - Write and delete operations update or **invalidate** the cache as needed.

- **Why this approach?**
    - Keeps caching logic separate from business logic.
    - Easy to test and extend.
    - Uses in-memory storage for simplicity—no external cache is required.

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
