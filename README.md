# Flashat

A real-time chat application built with Go, React, PostgreSQL, Redis, RabbitMQ, and WebSocket.

🔗 **Live Demo:** https://flashatapp.com

---

## Demo Accounts

Two accounts are pre-configured and already friends with each other so you can try the app immediately.

| Account | Email | Password |
|---------|-------|----------|
| User 1  | one@one.com | 123456 |
| User 2  | two@two.com | 123456 |

---

## Features

**Messaging**
- Send and receive messages in real time via WebSocket
- Messages show sending, sent, and failed states
- Unread message count per conversation
- Read receipts

**Direct Conversations**
- Start a direct conversation with any friend
- Messages sync across multiple tabs and devices

**Group Conversations**
- Create a group with multiple participants
- Leave a group at any time
- Creator role is automatically transferred if the creator leaves

**Friends**
- Send a friend request by email
- Accept, decline, or cancel pending requests
- Remove a friend
- Block a user

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go, Gin |
| Frontend | React, TypeScript |
| Database | PostgreSQL |
| Session | Redis |
| Message Queue | RabbitMQ |
| Real-time | WebSocket |
| Reverse Proxy | Nginx |
| Deployment | Docker Compose, AWS EC2 |
| CI/CD | GitHub Actions |