# go-auth-service

A standalone Go authentication service that any app or backend can use for sign-up, login, and token management.

This project is designed to demonstrate:
- A Go + Fiber + PostgreSQL + Redis service running fully in Docker
- Short-lived JWT access tokens (HS256) with rotating refresh tokens
- Cookie or Bearer token auth, so web apps, mobile apps, and other backends can all use it
- Layered code: routes, handlers, services, repositories
- Brute-force protection and login history

---

## ✨ Features

- Sign up and sign in with email and password.
- Get a JWT access token and a refresh token after sign up, sign in, or refresh.
- Send the token as an `httpOnly` cookie or as an `Authorization: Bearer` header.
- Refresh tokens rotate: each one works only once.
- An email is locked for a while after too many failed logins.
- Every register, login, logout, failed, and locked attempt is saved with IP and user agent.

---

## 🧠 Business Rules

1. The service stores identity only: id, email, and password hash. Names, avatars, and roles belong in your own app, keyed by the user ID (`sub` in the token).
2. Emails are trimmed and lowercased, so `Test@Example.com` and `test@example.com` are the same user.
3. Passwords must be 8–72 characters and are hashed with bcrypt.
4. After `LOGIN_MAX_ATTEMPTS` failed logins, the email is locked for `LOGIN_LOCK_DURATION` and gets `429` with a `Retry-After` header. A successful login resets the counter.
5. `/refresh` revokes the old refresh token and returns a new pair.
6. `/logout` revokes the refresh token. The access token stays valid until it expires.
7. A user with `is_active = false` cannot sign in, refresh, or read their profile (`403`).

---

## ⚙️ Configuration

Copy the example file, then edit `.env`:

```bash
cp .env.example .env
```

| Variable | Description |
|---|---|
| `APP_NAME` | Project name, also the default `JWT_ISSUER` |
| `APP_ENV` | `dev` or `prod`. In `prod`, `JWT_SECRET` is required and cookies are `Secure` by default. |
| `PG_HOST`, `PG_PORT`, `PG_USER`, `PG_PASSWORD`, `PG_DATABASE`, `PG_SCHEMA`, `PG_SSL_MODE` | PostgreSQL connection |
| `REDIS_ADDRESS`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB` | Redis connection |
| `SERVER_PORT` | API port (default `3001`) |
| `CORS_ORIGINS` | Allowed origins, separated by `,` |
| `PROXY_HEADER` | Header that holds the real client IP behind a proxy, such as `X-Forwarded-For` |
| `JWT_SECRET` | Token signing key, at least 32 characters. Generate one with `openssl rand -hex 48`. |
| `JWT_ISSUER`, `JWT_AUDIENCE` | `iss` and `aud` in the token |
| `ACCESS_TOKEN_TTL`, `REFRESH_TOKEN_TTL` | Token lifetimes (default `15m` and `720h`) |
| `LOGIN_MAX_ATTEMPTS`, `LOGIN_LOCK_DURATION` | Failed logins allowed before the email is locked, and for how long |
| `COOKIE_ENABLED` | Set token cookies on responses. Use `false` if no client is a browser. |
| `COOKIE_NAME`, `COOKIE_REFRESH_NAME` | Cookie names |
| `COOKIE_DOMAIN` | Use `.example.com` to share cookies across subdomains |
| `COOKIE_SAME_SITE`, `COOKIE_SECURE` | `lax`, `strict`, or `none`. `none` requires HTTPS. |

If `JWT_SECRET` is empty in `dev`, a temporary secret is used and tokens stop working after a restart. Changing the secret invalidates every issued token.

---

## 🚀 How to Run

### 1. Start the containers
```bash
docker compose up -d --build
```

The `go-auth-service` container stays idle until you start the app. Tables are created automatically on start.

### 2. Start the API
```bash
docker exec -it go-auth-service go run .
```

- API: http://localhost:3001/health

### Development notes
- To run Go on your machine instead (for example, to debug in an IDE), start only the databases. `.env` already points at `localhost`.
  ```bash
  docker compose up -d go-auth-db redis
  go run .
  ```
- Run the tests:
  ```bash
  go test ./...
  ```
- Build the production image, a static binary on debian-slim that runs as a non-root user:
  ```bash
  docker build --target prod -t go-auth-service .
  ```

---

## 🔌 API

All paths start with `/api/v1/auth`. Every response uses the shape `{ "status": "ok" | "error", "message": "...", "data": ... }`.

| Method | Path | Access | Description |
|---|---|---|---|
| `GET` | `/health` | Anyone | Health check (no `/api/v1/auth` prefix) |
| `POST` | `/register` | Anyone | Sign up (`email`, `password`) |
| `POST` | `/login` | Anyone | Sign in (`email`, `password`) |
| `POST` | `/refresh` | Anyone | New tokens (`refresh_token`, or the refresh cookie) |
| `POST` | `/logout` | Anyone | Revoke the refresh token (`refresh_token`, or the refresh cookie) |
| `GET` | `/profile` | User | Current user |

A user request needs `Authorization: Bearer <access_token>` or the access cookie.

| Status | Meaning |
|---|---|
| `400` | Validation failed |
| `401` | Wrong credentials, or a bad or expired token |
| `403` | User is inactive |
| `409` | Email already registered |
| `429` | Too many failed logins |

### Using the token in another backend

Give the backend the same `JWT_SECRET` and verify the token with any JWT library: algorithm `HS256`, `iss` equal to `JWT_ISSUER`, and `exp` not passed. Use `sub` as the user ID.

If a backend should not hold the secret, call `GET /api/v1/auth/profile` with the user's token instead.

---

## 📊 Example Output

`POST /api/v1/auth/login`:

```json
{
  "status": "ok",
  "message": "Login successful.",
  "data": {
    "access_token": "eyJ...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "opaque-string",
    "refresh_expires_in": 2592000,
    "user": { "id": "aa9e365e-1c51-45d9-b85d-cf2d82a3dcc8", "email": "test@example.com" }
  }
}
```

---

## 📁 Project Structure

```
go-auth-service/
├── bootstrap/          # Connects config, databases, repositories, and services
├── config/             # Env, JWT secret, Redis client
├── db/                 # Postgres connection and migration
├── dto/                # Request and response bodies
├── enums/              # App environments and login history actions
├── handler/            # Read the request, call a service, send the response, set cookies
├── middleware/         # Auth and CORS
├── models/             # Users and login history tables
├── repository/         # Postgres and Redis queries
├── routes/             # Paths and middleware
├── service/            # Business logic: auth, lockout, tokens
├── utils/              # Response helpers and validation
├── main.go             # Start the server
├── docker-compose.yml
├── Dockerfile
└── .env.example
```

---

## 📄 License
MIT License. See [LICENSE](LICENSE).
