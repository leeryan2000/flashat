# API Reference

Every route is under `/api`. All routes except `POST /login`, `POST /user/register`, and `GET /status` need the `session_id` cookie from login. For Postman-style testing, `?session_id=<id>` works too.

Errors come back as `{"error": "<message>"}` with a 4xx/5xx status.

## Health

| Method | Path | Notes |
|---|---|---|
| GET | `/status` | `{"message":"Ok"}` |

## Auth

| Method | Path | Body | Notes |
|---|---|---|---|
| POST | `/login` | `{email, password}` | Sets the `session_id` cookie and returns the user |
| DELETE | `/auth/logout` | – | Deletes the session and clears the cookie |
| GET | `/auth/me` | – | Current user |

## Users

| Method | Path | Body / Query | Notes |
|---|---|---|---|
| POST | `/user/register` | `{email, password, name, code}` | `code` must match `REGISTER_CODE` |
| GET | `/user/auth/all` | – | All users |
| GET | `/user/auth/:id` | – | User by UID |
| PUT | `/user/auth/name` | `{name}` | Rename the current user |
| GET | `/user/auth/avatar/upload-url` | – | Presigned S3 PUT URL (5 min, `image/jpeg`) |
| PUT | `/user/auth/avatar` | `{avatar_url}` | Save the avatar URL after uploading |
| DELETE | `/user/auth/avatar` | – | Delete the S3 object and clear the URL |

## Conversations

| Method | Path | Body / Query | Notes |
|---|---|---|---|
| POST | `/conversation/group` | `{name, participantIds: [uid]}` | Caller becomes `creator` |
| POST | `/conversation/direct` | `{target_uid}` | |
| GET | `/conversation/` | – | Conversations for the current user |
| GET | `/conversation/summary` | – | Sidebar summaries: last message, unread count, participants |
| GET | `/conversation/:conversation_id` | – | |
| GET | `/conversation/:conversation_id/participant` | – | |
| POST | `/conversation/participant` | `{conversation_id, participant_uid, role}` | Add a member |
| PUT | `/conversation/participant` | `{conversation_id, participant_uid, role}` | Change a member's role |
| DELETE | `/conversation/participant` | `{conversation_id, participant_uid}` | Remove a member |
| PUT | `/conversation/read_seq` | `{conversation_id, last_read_seq}` | |
| GET | `/conversation/read_seq` | `?conversation_id=` | Current user's last read seq |
| DELETE | `/conversation/:conversation_id/leave` | – | Leave a group. Passes `creator` on if needed |

## Messages

| Method | Path | Query | Notes |
|---|---|---|---|
| GET | `/message/latest/:conversation_id` | `limit` (default 50) | Newest messages |
| GET | `/message/before/:conversation_id` | `seq`, `limit` (default 50) | Older page, for scrollback |

## Friendships

| Method | Path | Body | Notes |
|---|---|---|---|
| GET | `/friendship` | – | Friends and pending requests |
| POST | `/friendship/request` | `{email}` | |
| POST | `/friendship/accept` | `{uid}` (requester) | Creates a direct conversation |
| DELETE | `/friendship/delete/:friend_uid` | – | Removes the friend and the direct conversation |
| DELETE | `/friendship/decline/:friend_uid` | – | |
| DELETE | `/friendship/cancel/:friend_uid` | – | Cancel a request you sent |
| POST | `/friendship/block/:uid` | – | |

## WebSocket

`GET /api/websocket/ws` (`ws://localhost:8080/api/websocket/ws` in development). Authentication works the same as for REST. Each user can have at most 5 connections; past that the server answers `429`.

Every frame is a JSON `wire.Msg`:

```json
{
  "type": "chat | ack | read | ...",
  "conversation_id": "uuid",
  "client_msg_id": "uuid",
  "server_msg_id": "uuid",
  "from_uid": "uuid",
  "from_name": "string",
  "seq": 42,
  "ts": 1730000000000,
  "body": { "text": "hello" },
  "last_read_seq": 42
}
```

### Client → server

| type | Fields | Effect |
|---|---|---|
| `chat` | `conversation_id`, `client_msg_id`, `from_uid`, `body` | Queued, saved, then acked and broadcast |
| `read` | `conversation_id`, `from_uid`, `last_read_seq` | Updates the read marker |

The server ignores any other `type`. A malformed frame closes the connection.

The server always replaces `from_uid` with the UID of the authenticated session, so clients can't send messages or read receipts as someone else. Clients may still include it.

### Server → client

| type | Sent to | Fields |
|---|---|---|
| `ack` | Sender's connection(s) | `client_msg_id`, `server_msg_id`, `seq`, `ts`, `from_name`, `body` |
| `chat` | Other participants and the sender's other tabs | `server_msg_id`, `seq`, `ts`, `from_uid`, `from_name`, `body` |

`ts` is in epoch milliseconds.
