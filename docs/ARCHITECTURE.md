# Architecture

## System overview

```
Browser (React SPA)
   │  HTTPS (REST /api/*)      WSS (/api/websocket/ws)
   ▼
Nginx ── serves frontend/dist, terminates TLS (Cloudflare origin cert)
   │
   ▼
Go backend (Gin, :8080)
   ├── PostgreSQL   users, conversations, messages, friendships
   ├── Redis        login sessions (session:<id> → uid, 3-day TTL)
   ├── RabbitMQ     durable queue "chat_messages"
   └── AWS S3       avatar images (presigned PUT uploads)

Promtail → Loki → Grafana   container log collection and dashboards
```

All services run under one `docker-compose.yml` on an AWS EC2 host. In production the frontend and API share an origin, so cookies work without cross-site settings.

## Backend layout

| Package | Responsibility |
|---|---|
| `main.go` | Logging setup, HTTP server, graceful shutdown (SIGINT/SIGTERM) |
| `config/` | Load `.env` / process env into `Configuration` |
| `server/` | `StartServer` dependency wiring; `Hub` and `Client` for WebSocket connections |
| `transport/`, `transport/routes/` | Router, CORS, route groups, auth middleware attachment |
| `middleware/` | `Authenticate` — session cookie → UID |
| `handlers/` | Gin handlers for REST and the WebSocket upgrade |
| `service/` | Chat pipeline: `MessageService`, `MessageWorker`, observers |
| `repo/`, `repo/impl/` | Repository interfaces and pgx implementations |
| `models/` | DB-shaped structs |
| `wire/` | Client-facing JSON types (`wire.Msg`, summaries) |
| `db/` | Postgres pool + embedded migrations, Redis, RabbitMQ, S3 clients |
| `utils/` | bcrypt hashing, transaction helpers, OTP generator |
| `mocks/` | testify mocks for the `service` interfaces |

## Startup

`server.StartServer`:

1. Loads config and runs embedded migrations (`db/migrations/*.sql`) with golang-migrate.
2. Opens the pgx pool, Redis client, and RabbitMQ channel (declares `chat_messages`).
3. Creates the `Hub`, repos, and S3 clients (default AWS credential chain — instance role in prod).
4. Builds `MessageService`, registers `AckObserver` and `BroadcastObserver`.
5. Starts `MessageWorker` and a 30-second task that logs connection counts.

On shutdown, `main.go` stops the HTTP server (10s), waits for the worker pool to drain (`ShutdownTimeout` = 13s), then closes RabbitMQ, Redis, and Postgres.

## Real-time messaging

### Connections

- `GET /api/websocket/ws` is authenticated by the session cookie, then upgraded.
- Each connection becomes a `Client` with a `ReadPump` and `WritePump` goroutine and a 256-message send buffer.
- The `Hub` indexes clients by UID, so a user can have several tabs or devices. The limit is 5 connections per user, checked atomically inside `Hub.Register`.
- Keepalive: the server pings every 25s and expects a pong within 60s.
- In production, the WebSocket origin must be `https://flashatapp.com` or `https://www.flashatapp.com`.

### Sending a message

```
Client ──chat──▶ ReadPump ─▶ MessageService.handleChat
                               │  validate IDs, assign server_msg_id
                               ▼
                         RabbitMQ "chat_messages" (persistent)
                               │
                               ▼
                   MessageWorker dispatcher (single consumer)
                               │  fnv32a(conversation_id) % 4
                               ▼
                       worker[i] ─▶ MessageService.ProcessChat
                                      1. SaveMessage (assigns seq)
                                      2. ListParticipantByID
                                      3. observers
                                         ├─ AckObserver  → sender: type "ack"
                                         └─ BroadcastObserver → other participants + sender's other tabs: type "chat"
```

- **Ordering:** each conversation always hashes to the same worker, so its messages are saved in order. Different conversations run in parallel.
- **Sequence numbers:** `utils.SaveMessageWithTx` increments `conversation_counters.last_seq` and inserts the message in one statement. The insert only happens if the sender is a participant.
- **Retries:** `ProcessChat` gets 3 attempts, 200ms apart, inside an 8s timeout. After that the delivery is nacked without requeue (dropped).
- **Client status:** the frontend marks a message `sending`, flips it to `sent` on `ack` (matched by `client_msg_id`), or `failed` after 10s with no ack. The backend 8s timeout is set below the 10s frontend timeout on purpose.

### Read receipts

The client sends `{"type":"read", "conversation_id", "last_read_seq"}`. `handleRead` writes it straight to `conversation_participants.last_read_seq` without going through the queue. Unread count = `last_seq - last_read_seq`.

## Data model

| Table | Key columns |
|---|---|
| `users` | `uid` PK, `email` unique, `hashed_password`, `user_avatar_url`, `role` (`user`/`admin`) |
| `conversations` | `id` PK, `type` (`direct`/`group`), `group_name`, `group_avatar_url` |
| `conversation_participants` | PK (`conversation_id`, `uid`), `role` (`creator`/`admin`/`member`), `last_read_seq` |
| `conversation_counters` | `conversation_id` PK, `last_seq` |
| `messages` | `id` PK, `conversation_id`, `seq` (unique per conversation), `from_uid`, `body` jsonb (`{"text": ...}`) |
| `friendships` | unique (`requester_uid`, `receiver_uid`), `status` (`pending`/`accepted`/`blocked`) |

Rules worth knowing:

- Accepting a friend request creates a direct conversation. Removing a friend deletes that conversation.
- When a group's creator leaves, the next member is promoted to `creator`.

## Authentication

1. `POST /api/login` checks the bcrypt hash, writes `session:<uuid>` → uid to Redis (3-day TTL), and sets an HttpOnly `session_id` cookie. In production the cookie is `Secure` on `.flashatapp.com`.
2. `middleware.Authenticate` resolves the cookie to a UID and stores it as `c.Set("uid", ...)`.
3. Registration requires `REGISTER_CODE` (a shared invite code until email verification exists).

## Frontend

- `src/api/api.ts`: `fetch` wrapper that sends cookies and throws on `{ error }` responses.
- `src/context/`
  - `AuthContext`: current user, login/logout, profile and avatar updates.
  - `WebSocketContext`: a single socket that reconnects every 5s unless closed on purpose.
  - `ChatContext`: conversation and message state, pending/ack handling, read receipts.
  - `ThemeContext`: light/dark mode.
- `src/wire/`: server DTO types plus `dtoTo*` converters to camelCase client types.
- `src/pages/`, `src/components/`, `src/layouts/`, `src/routes/`: UI and routing (`/`, `/login`, `/register`, `/friends`, `/profile`, `/settings`).

## Avatar upload

1. `GET /api/user/auth/avatar/upload-url` returns a presigned S3 PUT URL for `avatars/<uid>.jpg` (5 min, `image/jpeg`).
2. The browser uploads straight to S3.
3. `PUT /api/user/auth/avatar` saves the public URL on the user. `DELETE` removes the S3 object and clears the URL.

## Observability

The backend logs JSON through `log/slog` to stdout. Promtail reads Docker logs and sends them to Loki. Grafana listens on `127.0.0.1:3000` on the EC2 host; reach it over an SSH tunnel.
