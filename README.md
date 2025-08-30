# 💸 Paysplit Backend (🇮🇳 UPI Ready)

Backend service for **Paysplit**, a modern Splitwise-style bill splitting app designed for India.  
Built with **Go (Fiber)**, **Postgres**, and **Redis**, with support for **UPI deep links** for instant settlement.

---

## ✨ Features
- 👤 User management (create, update, delete)
- 👥 Groups & members (create groups, add members)
- 💰 Expenses (equal, exact, shares, percent)
- 📊 Balances & simplify debts (minimal transfers)
- 🔗 Generate UPI deep links (`upi://` + `paysplit://`)
- 🔒 JWT authentication (planned)
- 📈 Structured logging with Zap
- 🐳 Dockerized local setup (Postgres + Redis)

---

## 🏗️ Tech Stack
- **Go** (Fiber web framework)
- **Postgres** (SQL database)
- **Redis** (caching, reminders)
- **Zap** (structured logging)
- **Docker Compose** (local infra)

---

## 🚀 Getting Started

### Prerequisites
- Go `1.21+`
- Docker + Docker Compose
- Postgres client (`psql`)

### Run locally
```bash
# start Postgres + Redis
docker compose up -d db redis

# apply migrations
psql "postgres://paysplit:paysplit@localhost:5432/paysplit?sslmode=disable" -f db/init.sql

# run API
go run ./cmd/api
