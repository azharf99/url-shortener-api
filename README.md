# URL Shortener API

A robust, secure, and scalable URL Shortener service built with Go, following Clean Architecture principles.

## 🚀 Features

- **URL Shortening:** Generate unique, short aliases for long URLs.
- **Redirection:** Fast redirection from short codes to original destination.
- **User Management:**
    - Secure registration and login system.
    - Google OAuth 2.0 integration.
    - Admin dashboard for user CRUD operations.
- **Security & Reliability:**
    - **JWT Authentication:** Secure access to protected resources.
    - **Rate Limiting:** Protect the API from abuse using Redis-backed rate limiting.
    - **ReCAPTCHA Integration:** Enhanced security for sensitive actions (Login, Register, URL creation).
    - **CORS Support:** Configurable allowed origins for web applications.
- **Analytics:** Basic click tracking for each shortened URL.
- **Logging:** Structured logging using `uber-go/zap`.
- **Docker Ready:** Fully containerized for easy deployment.

## 🛠️ Tech Stack

- **Language:** [Go](https://golang.org/) (v1.26+)
- **Framework:** [Gin Web Framework](https://gin-gonic.com/)
- **Database:** [PostgreSQL](https://www.postgresql.org/) with [GORM](https://gorm.io/)
- **Cache & Rate Limiting:** [Redis](https://redis.io/)
- **Authentication:** [JWT](https://github.com/golang-jwt/jwt)
- **Logging:** [Zap](https://github.com/uber-go/zap)
- **Containerization:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

## 🏗️ Architecture

This project follows **Clean Architecture** (Repository, Usecase, Delivery) to ensure maintainability, testability, and separation of concerns.

```text
├── cmd/api             # Application entry point
├── internal/
│   ├── config          # Database, Redis, and Logger configurations
│   ├── delivery/http   # HTTP handlers and middlewares
│   ├── domain          # Domain models and interfaces
│   ├── repository      # Data access layer
│   ├── usecase         # Business logic layer
│   └── utils           # Helper functions (Auth, etc.)
```

## 🚦 Getting Started

### Prerequisites

- Go 1.26 or later
- Docker & Docker Compose
- Redis (optional for local run, included in Docker)
- PostgreSQL (optional for local run, included in Docker)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/azharf99/url-shortener-api.git
   cd url-shortener-api
   ```

2. **Configure Environment Variables:**
   Copy `.env.example` to `.env` and update the values.
   ```bash
   cp .env.example .env
   ```

3. **Run with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

4. **Run Locally:**
   Make sure you have PostgreSQL and Redis running, then:
   ```bash
   go run cmd/api/main.go
   ```

## 📖 API Documentation

### Public Endpoints
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/:shortCode` | Redirect to original URL |
| `POST` | `/api/register` | Register a new user (Captcha required) |
| `POST` | `/api/login` | Login and get JWT (Captcha required) |
| `POST` | `/api/google-login` | Login with Google OAuth (Captcha required) |

### Protected Endpoints (Requires JWT)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/shorten` | Create a new short URL (Captcha required) |
| `GET` | `/api/urls` | List user's shortened URLs |
| `PUT` | `/api/urls/:id` | Update a short URL (Captcha required) |
| `DELETE` | `/api/urls/:id` | Delete a short URL (Captcha required) |

### Admin Endpoints (Requires Admin Role)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/admin/users` | List all users |
| `POST` | `/api/admin/users` | Create a user by admin (Captcha required) |
| `GET` | `/api/admin/users/:id` | Get user details |
| `PUT` | `/api/admin/users/:id` | Update user (Captcha required) |
| `DELETE` | `/api/admin/users/:id` | Delete user (Captcha required) |

## 📄 License

This project is licensed under the **Apache License 2.0**.

**Copyright © 2026 Azhar Faturohman Ahidin**

Under the terms of this license, any use, reproduction, or distribution of this work must include the above copyright notice and mention the author **Azhar Faturohman Ahidin**.

---
Made with ❤️ by [Azhar Faturohman Ahidin](https://github.com/azharf99)
