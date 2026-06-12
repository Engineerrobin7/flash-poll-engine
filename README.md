# ⚡ Flash-Poll Engine

> **Internal real-time polling for teams that move fast.**

Built for the LeMiCi Engineering Technical Assessment, Flash-Poll is a high-velocity decision engine. It focuses on **atomic integrity**, **real-time synchronization**, and a **no-nonsense user experience**.

---

## 🛠 The Stack

- **Backend**: Go 1.25 (Standard library + `chi` router)
- **Database**: SQLite (Pure Go / CGO-free via `modernc.org/sqlite`)
- **Frontend**: React 19 + Vite (Zero-dependency Neo-Brutalist CSS)
- **Real-time**: Server-Sent Events (SSE) for multi-client live sync
- **Deployment**: Docker & Docker Compose support

---

## 🏗 Why this Architecture?

### 1. Atomic Voting (No Race Conditions)
Most polling apps fail under concurrent load because they read a value, increment it in memory, and save it back. Flash-Poll uses **direct SQL atomic increments**:
`UPDATE options SET vote_count = vote_count + 1 WHERE id = ?`.
This ensures that even if 100 people vote at the exact same millisecond, every single vote is counted accurately.

### 2. Live "Signal" Broadcasting
Instead of heavy WebSockets, I implemented a custom **SSE Broker** in Go. When a vote is cast, the backend broadcasts the update to all connected clients. The UI updates live without the user ever hitting refresh.

### 3. Neo-Brutalist UX
No rounded corners, no soft gradients. The UI is designed to feel like a high-contrast engineering tool. High-vis borders, tactile button feedback, and themed category accents for **Tech**, **Business**, and **Design**.

---

## 🚦 Getting Started

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

## 📡 API Contract

| Endpoint | Method | Action |
| :--- | :--- | :--- |
| `/api/polls` | `GET` | Fetch the live feed |
| `/api/polls` | `POST` | Create a new signal |
| `/api/polls/:id/vote` | `PATCH` | Record an atomic vote |
| `/api/polls/:id` | `DELETE` | Purge poll & associated data |
| `/api/events` | `GET` | Real-time SSE stream |

---

## 🧪 Robustness Testing
The system includes a `stress_test.go` script that blasts the server with 100 concurrent votes.
- **Panic Recovery**: Middleware ensures the binary doesn't crash on invalid inputs.
- **Deadlock Protection**: Non-blocking SSE broker prevents slow clients from hanging the server.

---
*Built with focus and coffee for the LeMiCi Technical Assessment.*
