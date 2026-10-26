# Development

## Prerequisites

- Docker Desktop
- Node.js 20+

Go is not required: the backend is built inside Docker.

## Start the stack

The backend and its infrastructure (PostgreSQL, Redis, RabbitMQ) run in Docker through `docker-compose.dev.yml`. The frontend runs on the host with Vite.

```bash
docker compose -f docker-compose.dev.yml up --build --watch
```

`--watch` rebuilds and restarts only the backend when a file under `backend/` changes; the other services keep running. Drop it to rebuild manually instead (`up --build`).

In a second terminal:

```bash
cd frontend
npm install
npm run dev
```

| Service | URL |
|---------|-----|
| Frontend | http://localhost:5173 |
| Backend | http://localhost:8080/api/status |
| RabbitMQ UI | http://localhost:15672 (`flashat` / `flashat`) |

## Database

Migrations in `backend/db/migrations/` are embedded in the binary and applied on every backend start. A new migration file is picked up on the next rebuild.

The local database starts empty. Register through the app; the verification code is `123456` (`REGISTER_CODE`).

## Common commands

All commands need `-f docker-compose.dev.yml`; without it Compose uses the production `docker-compose.yml`.

| Command | What it does |
|---------|--------------|
| `docker compose -f docker-compose.dev.yml logs -f backend` | Follow backend logs |
| `docker compose -f docker-compose.dev.yml up -d --build backend` | Rebuild only the backend |
| `docker compose -f docker-compose.dev.yml down` | Stop everything, keep data |
| `docker compose -f docker-compose.dev.yml down -v` | Stop everything and delete the database and queue data |

## Configuration

Values are read from the shell or a root `.env` file (gitignored).

| Variable | Default | Notes |
|----------|---------|-------|
| `BACKEND_PORT` | `8080` | If changed, update `VITE_API_URL` and `VITE_WS_URL` in `frontend/.env.development` |
| `DB_PORT`, `REDIS_PORT`, `MQ_PORT`, `MQ_UI_PORT` | `5432`, `6379`, `5672`, `15672` | Host ports |
| `REGISTER_CODE` | `123456` | Verification code for registration |
| `AWS_REGION` | `ap-southeast-2` | |
| `AWS_S3_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` | empty | Needed only for avatar uploads |

Avatar uploads need an IAM user with write access to the bucket's `avatars/` prefix. In production the backend uses the EC2 instance role instead, so no keys are stored there.
