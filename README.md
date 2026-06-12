#  Flash-Poll Engine

> **Internal real-time polling for teams that move fast.**

Built for the LeMiCi Engineering Technical Assessment, Flash-Poll is a high-velocity decision engine. It focuses on **atomic integrity**, **real-time synchronization**, and a **no-nonsense user experience**.

---

##  The Stack

- **Backend**: Go 1.25 (Standard library + `chi` router)
- **Database**: SQLite (Pure Go / CGO-free via `modernc.org/sqlite`)
- **Frontend**: React 19 + Vite (Zero-dependency Neo-Brutalist CSS)
- **Real-time**: Server-Sent Events (SSE) for multi-client live sync
- **Deployment**: Docker & Docker Compose support

---

##  Why this Architecture?

### 1. Atomic Voting (No Race Conditions)
Most polling apps fail under concurrent load because they read a value, increment it in memory, and save it back. Flash-Poll uses **direct SQL atomic increments**:
`UPDATE options SET vote_count = vote_count + 1 WHERE id = ?`.
This ensures that even if 100 people vote at the exact same millisecond, every single vote is counted accurately.

### 2. Live "Signal" Broadcasting
Instead of heavy WebSockets, I implemented a custom **SSE Broker** in Go. When a vote is cast, the backend broadcasts the update to all connected clients. The UI updates live without the user ever hitting refresh.

### 3. Neo-Brutalist UX
No rounded corners, no soft gradients. The UI is designed to feel like a high-contrast engineering tool. High-vis borders, tactile button feedback, and themed category accents for **Tech**, **Business**, and **Design**.

---

##  Getting Started

### Prerequisites
- Go installed
- Node.js (v18+)
- (Optional) Docker

### 1. Fire up the Backend
```bash
cd backend
go mod tidy
go build -o server.exe cmd/server/main.go
./server.exe
```
*Server starts on `localhost:8080`. Database `flashpoll.db` is auto-initialized on first run.*

### 2. Launch the Frontend
```bash
cd frontend
npm install
npm run dev
```
*Interface available at `http://localhost:5173`.*

---

##  API Contract

| Endpoint | Method | Action |
| :--- | :--- | :--- |
| `/api/polls` | `GET` | Fetch the live feed |
| `/api/polls` | `POST` | Create a new signal |
| `/api/polls/:id/vote` | `PATCH` | Record an atomic vote |
| `/api/polls/:id` | `DELETE` | Purge poll & associated data |
| `/api/events` | `GET` | Real-time SSE stream |

---

##  Robustness Testing
The system includes a `stress_test.go` script that blasts the server with 100 concurrent votes.
- **Panic Recovery**: Middleware ensures the binary doesn't crash on invalid inputs.
- **Deadlock Protection**: Non-blocking SSE broker prevents slow clients from hanging the server.
- **Rate Limiting**: Built-in protection against automated vote-spamming (500ms cooldown per IP).
- **Graceful Shutdown**: Handles OS signals (SIGINT/SIGTERM) to close DB connections and finish in-flight requests cleanly.

## 🚀 Production Readiness
- **Observability**: Chi standard middleware for `RequestID` and `RealIP` tracking.
- **Health Checks**: Dedicated `/health` endpoint for uptime monitoring.
- **Accessibility**: ARIA-labeled components for screen-reader compatibility.
- **UX Scalability**: Category-based dashboard filtering for large datasets.

## ⚖️ Engineering Tradeoffs & Roadmap
1. **SSE vs WebSockets**: I chose SSE (Server-Sent Events) because the application is primarily read-heavy. SSE is more efficient, handles reconnection automatically (implemented with exponential backoff logic), and is simpler to scale than full-duplex WebSockets.
2. **SQLite vs PostgreSQL**: SQLite was chosen for this MVP to provide a "zero-config" experience. In a production environment with high write-concurrency, the repository layer is designed to be easily swapped for **PostgreSQL** to handle row-level locking more efficiently.
3. **Observability at Scale**: The current `Broker` uses in-memory channels. To scale this across multiple server instances (horizontal scaling), I would migrate the broadcasting logic to **Redis Pub/Sub**.
4. **Session Tracking**: To keep this MVP atomic and focused on performance, per-user session tracking (preventing multiple votes from the same user) was omitted but is the first item on the v2 roadmap.
5. **Rate Limiting**: I implemented a custom IP-based rate limiter that strips port information to prevent bypasses, ensuring the system remains protected against automated "vote-bombing."

---
*Built with focus and coffee for the LeMiCi Technical Assessment.*
