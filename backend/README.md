# ⚙️ Flash-Poll Backend (Go)

This is the core engine for the Flash-Poll system. It's a modular Go service designed for concurrency and data safety.

## Highlights

- **CGO-Free SQLite**: Uses `modernc.org/sqlite` for maximum portability without needing a C compiler.
- **Atomic Integrity**: Zero "read-modify-write" cycles. All increments happen at the database level.
- **Real-time Broker**: Custom SSE implementation for low-latency broadcasting.
- **Clean Architecture**: 
  - `domain/`: Models and Enums
  - `repository/`: SQL & Transactions
  - `service/`: Business Logic & Validation
  - `handlers/`: HTTP Interface
  - `middleware/`: Logging & Panic Recovery

## Running Locally

```bash
go mod tidy
go run cmd/server/main.go
```

## Running the Stress Test

To prove the system can handle concurrent loads without losing votes or crashing:

```bash
go run stress_test.go
```

---
*Status: Ready for Production Assessment.*
