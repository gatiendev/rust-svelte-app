# Complete Project Template – Lightweight Go + Svelte + Postgres Stack

A resource‑efficient, production‑ready template for building modern web applications. It combines a **Go (Golang) API**, a **Svelte frontend**, and a **PostgreSQL** database, all orchestrated with Docker Compose. Designed to be minimal, fast, and adaptable – perfect for startups or any project where efficiency matters.

---

## Overview

This template provides a full‑stack foundation with:

- **Go API** – High‑performance backend with JWT authentication, health checks, and database connectivity.
- **Svelte + Vite Frontend** – Real‑time capable UI (crypto ticker example included) that compiles to tiny static files.
- **PostgreSQL** – Relational database for persistent storage.
- **Docker Compose** – Unified development and production environment with minimal overhead.

The entire stack is optimized for low resource consumption, fast startup, and easy customization.

---

## Features

- **Modular Services** – Each component runs in its own container, following best practices.
- **Authentication Ready** – Go API includes JWT access/refresh token logic (configurable durations).
- **Real‑time WebSocket** – Frontend connects to the Go API’s WebSocket endpoint for live data.
- **Health Checks** – All services have health checks for orchestration reliability.
- **Hot Reload in Development** – Volume mounts enable live code updates without rebuilding.
- **Multi‑stage Docker Builds** – Production images are tiny (e.g., frontend ~23 MB, Go binary stripped).
- **Environment Variable Configuration** – Easy to adapt to different deployments.
- **PostgreSQL Persistence** – Database data stored in a Docker volume.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Docker Host                          │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │   Frontend   │◄──►│    Go API    │◄──►│  PostgreSQL  │  │
│  │  (Svelte)    │    │   (Golang)   │    │              │  │
│  │  port 5173   │    │   port 8000  │    │   port 5432  │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│         ▲                  ▲                    ▲           │
│         └──────────┬────────┴────────────────────┘           │
│                    │                                          │
│              app-network (bridge)                             │
└─────────────────────────────────────────────────────────────┘
```

- **Frontend** serves a single‑page application (SPA) via Nginx. It connects to the Go API over WebSocket (`ws://go_api:8000/ws`).
- **Go API** provides REST endpoints (e.g., `/health`) and a WebSocket feed (e.g., simulated price updates). It connects to PostgreSQL for user/auth data.
- **PostgreSQL** stores user credentials, session info, or any other relational data.

All services communicate through a dedicated bridge network (`app-network`).

---

## Efficiency and Resource Optimization

This template is built with resource efficiency as a core principle:

| Component       | Optimization Strategies                                                                                                                                                                                                                           |
|-----------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Go API**      | Compiled to a single binary; no interpreter overhead. Uses Gin in release mode. Alpine‑based Docker image keeps size small.                                                                                                                     |
| **Frontend**    | Svelte compiles away the framework, leaving minimal JavaScript. Vite tree‑shakes and minifies. Multi‑stage Docker build results in a static Nginx image (~23 MB).                                                                               |
| **PostgreSQL**  | Alpine‑based image (postgres:18-alpine) reduces disk and memory footprint. Health checks ensure quick recovery.                                                                                                                                  |
| **Docker Compose** | Lightweight networking, volume management, and only necessary services are started. Development mounts avoid duplication.                                                                                                                      |
| **Overall**     | No heavy runtimes (Node.js is only in build stage). All production images are based on Alpine Linux. Resource limits can be added easily in compose.                                                                                             |

---

## Project Structure

```
.
├── docker-compose.yml          # Main orchestration file
├── go_api/                     # Go backend service
│   ├── Dockerfile.dev           # Development build (with air for hot reload)
│   ├── Dockerfile.prod          # Production multi‑stage build
│   ├── main.go                  # Entry point (Gin server + WebSocket)
│   ├── handlers/                # HTTP and WebSocket handlers
│   ├── middleware/              # JWT auth, logging, etc.
│   ├── models/                  # Data models
│   ├── migrations/              # SQL migration files
│   └── go.mod / go.sum
├── frontend/                    # Svelte frontend
│   ├── Dockerfile.prod           # Production multi‑stage build (nginx)
│   ├── Dockerfile.dev            # Optional dev container (Node.js)
│   ├── src/
│   │   ├── App.svelte
│   │   ├── main.js
│   │   └── app.css
│   ├── public/
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── .env.example                 # Example environment variables
└── README.md                    # This file
```

---

## Getting Started

### Prerequisites

- **Docker** and **Docker Compose** (v2.0+)
- **Make** (optional, for convenience commands)

### Quick Start

1. **Clone the repository**  

   ```bash
   git clone <repository-url>
   cd <project-folder>
   ```

2. **Set environment variables**  
   Copy `.env.example` to `.env` and adjust values (especially JWT secrets).  

   ```bash
   cp .env.example .env
   ```

3. **Start all services**  

   ```bash
   docker-compose up -d
   ```

4. **Access the application**  
   - Frontend: http://localhost:5173  
   - Go API health check: http://localhost:8000/health  
   - Database: `localhost:5432` (credentials from `.env`)

5. **View logs**  

   ```bash
   docker-compose logs -f
   ```

6. **Stop services**  

   ```bash
   docker-compose down
   ```

---

## Service Configuration

### Go API

- **Environment variables** (set in `docker-compose.yml` or `.env`):  

  | Variable                | Description                          | Default           |
  |-------------------------|--------------------------------------|-------------------|
  | `DB_HOST`               | PostgreSQL host                      | `postgres`        |
  | `DB_PORT`               | PostgreSQL port                      | `5432`            |
  | `DB_USER`               | Database user                        | `auth_user`       |
  | `DB_PASSWORD`           | Database password                    | `auth_pass`       |
  | `DB_NAME`               | Database name                        | `auth_db`         |
  | `JWT_ACCESS_SECRET`     | Secret for access tokens             | `your-access...`  |
  | `JWT_REFRESH_SECRET`    | Secret for refresh tokens            | `your-refresh...` |
  | `ACCESS_TOKEN_DURATION` | Access token lifetime (Go duration)  | `15m`             |
  | `REFRESH_TOKEN_DURATION`| Refresh token lifetime               | `168h` (7 days)   |
  | `SERVER_PORT`           | Port the API listens on               | `8000`            |
  | `GIN_MODE`              | Gin mode (release/debug)             | `release`         |

- **Development**: The service uses `Dockerfile.dev` which installs `air` for hot reload. Source code is mounted so changes trigger rebuilds.
- **Production**: Use `Dockerfile.prod` for a multi‑stage build resulting in a tiny binary.

### Frontend

- **Build arguments** (passed in `docker-compose.yml`):  

  | Argument       | Description                          | Default                     |
  |----------------|--------------------------------------|-----------------------------|
  | `VITE_WS_URL`  | WebSocket URL for the frontend       | `ws://go_api:8000/ws`       |

- **Development**: A separate `Dockerfile.dev` can be used (not shown in compose) if you prefer to run the Vite dev server inside a container. The current compose uses the production Dockerfile and mounts `node_modules` to avoid overwriting, but for live development you might want to use a dev container or run `npm run dev` locally.
- **Production**: The `Dockerfile.prod` builds the app and serves it with Nginx on port 80.

### PostgreSQL

- **Image**: `postgres:18-alpine`
- **Environment**:  
  - `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` (set in compose)
- **Data persistence**: Docker volume `pg_data` mounted at `/var/lib/postgresql/data`
- **Health check**: `pg_isready` ensures the database is ready before the API starts.

---

## Development Workflow

### Without Containers (Local)

- **Go API**: `cd go_api && go run main.go` (requires Postgres running separately)
- **Frontend**: `cd frontend && npm install && npm run dev`

### With Containers (Recommended)

- **Full stack**: `docker-compose up -d`
- **Rebuild a single service**: `docker-compose up -d --build go_api`
- **View logs**: `docker-compose logs -f frontend`
- **Execute commands inside a container**:  

  ```bash
  docker-compose exec go_api /bin/sh
  docker-compose exec postgres psql -U auth_user auth_db
  ```

---

## Adding a New Feature / Customizing

This template is intentionally minimal so you can extend it for your startup idea. Here are common modifications:

### 1. Add New API Endpoints

Edit `go_api/main.go` or create new handlers in `handlers/`. The Gin router makes it easy to add REST or WebSocket routes.

### 2. Change the Frontend Symbol or Data Source

- Modify `frontend/src/App.svelte` to change the default symbol or the WebSocket message format.
- If you need multiple symbols, consider passing a prop or reading from URL.

### 3. Add Authentication UI

- The Go API already includes JWT secrets; you can build login/register endpoints and connect them to Svelte forms.
- Store tokens in `localStorage` or cookies.

### 4. Use a Different Database

- Replace the `postgres` service with MySQL, MongoDB, etc. Adjust the Go API’s database driver accordingly.

### 5. Add Background Jobs

- Use Go routines or a separate worker service (add a new container to compose).

### 6. Set Resource Limits

To further optimize, add resource constraints in `docker-compose.yml`:

```yaml
services:
  go_api:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
  postgres:
    deploy:
      resources:
        limits:
          memory: 512M
```

---

## Production Deployment

### Building Production Images

You can build production images manually:

```bash
# Go API
docker build -f go_api/Dockerfile.prod -t myapp-api:latest ./go_api

# Frontend
docker build -f frontend/Dockerfile.prod --build-arg VITE_WS_URL=wss://api.example.com/ws -t myapp-frontend:latest ./frontend
```

Then push to a registry and deploy on any orchestrator (Docker Swarm, Kubernetes, or single server with docker run).

### Using Docker Compose in Production

For a single‑server production setup, you can use the same `docker-compose.yml` with a `.env` file containing production secrets. Ensure you:

- Remove volume mounts that expose source code.
- Use `restart: always` (already set).
- Consider adding a reverse proxy (like Traefik or Nginx) in front of the services.
- Enable logging drivers and resource limits.

---

## Environment Variables

Create a `.env` file in the project root (see `.env.example`):

```bash
# PostgreSQL
POSTGRES_USER=auth_user
POSTGRES_PASSWORD=strongpassword
POSTGRES_DB=auth_db

# Go API
JWT_ACCESS_SECRET=your-access-secret
JWT_REFRESH_SECRET=your-refresh-secret
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h

# Frontend (build-time)
VITE_WS_URL=ws://go_api:8000/ws   # in production, use wss://yourdomain.com/ws
```

> **Note**: The frontend `VITE_WS_URL` is a build argument, not a runtime variable. If you need to change it without rebuilding, consider using a runtime configuration script.

---

## Health Checks

All services define health checks to ensure orchestration reliability:

- **Go API**: `wget --spider http://localhost:8000/health`
- **Postgres**: `pg_isready -U auth_user -d postgres`
- **Frontend**: Nginx itself is healthy if it responds on port 80 (implicitly checked by Docker if `HEALTHCHECK` is added; you can add one if needed).

Docker Compose will wait for dependencies (`depends_on` with `condition: service_healthy` can be added for stricter ordering).

---

## Troubleshooting

### Frontend cannot connect to WebSocket

- Ensure `VITE_WS_URL` in the build args points to the correct service name (`go_api`) and port.
- Check that the Go API is running and the WebSocket endpoint is available (`/ws`).
- If accessing from a browser outside Docker, the WebSocket URL must be the host's address (e.g., `ws://localhost:8000/ws`). In production, use the public domain.

### Go API cannot connect to Postgres

- Verify that the `DB_HOST` is set to `postgres` (the service name) and the database credentials match.
- Check if Postgres is healthy: `docker-compose ps postgres`.
- Run `docker-compose logs postgres` for errors.

### Permission issues with mounted volumes

- On Linux, if you encounter permission errors with `go_mod_cache` volume, ensure the container user has write permissions. You can set user ID in the Dockerfile or run with `user: "${UID}:${GID}"` in compose.

---

## License

This project template is open source and available under the [MIT License](LICENSE). Feel free to use it for any personal or commercial project.

---

*Built with ❤️ for efficient startups – Go, Svelte, and Docker working together.*
