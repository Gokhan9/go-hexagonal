# Go Hexagonal Wallet & Auth System

This project is a monorepo containing a robust, thread-safe Wallet and Authentication API implemented using **Hexagonal Architecture** (Ports & Adapters) in Go, accompanied by a modern **React/Vite** frontend.

## Project Structure

```text
.
├── backend/    # Go Backend (Hexagonal Architecture)
└── frontend/   # React/TypeScript Frontend
```

## Backend (Go)

The backend follows the Hexagonal Architecture pattern to decouple core business logic from external technologies.

### Key Features
- **Concurrency Management:** Thread-safe operations using `sync.Mutex` and mechanisms to handle race conditions.
- **Idempotency:** Protects against duplicate requests using `X-Idempotency-Key` headers.
- **Transaction Management:** Ensures atomicity for financial operations (e.g., Transfers) using database transactions.
- **Authentication:** Secure user registration and login using JWT.
- **Rate Limiting:** Protects API endpoints from abuse.
- **API Documentation:** Interactive Swagger UI available.

### API Endpoints
All API endpoints are documented using Swagger/OpenAPI.
- **Base URL:** `http://localhost:8080`
- **Documentation:** `http://localhost:8080/swagger/index.html`

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| POST | `/wallets` | Create a new wallet |
| GET | `/wallets/{id}` | Get wallet details |
| POST | `/wallets/{id}/deposit` | Deposit funds |
| POST | `/wallets/{id}/withdraw` | Withdraw funds |
| POST | `/wallets/{id}/transfer` | Transfer funds between wallets |
| GET | `/wallets/{id}/balance` | Get wallet balance |
| GET | `/wallets/{id}/transactions` | Get transaction history |

### Getting Started (Backend)
Navigate to the `backend/` directory:
```bash
cd backend
go run cmd/main.go
```
The server will start on port `8080`.

## Frontend (React/Vite)

A React application built with Vite and TypeScript for interacting with the backend API.

### Getting Started (Frontend)
Navigate to the `frontend/` directory:
```bash
cd frontend
npm install
npm run dev
```

## Testing

### Backend
Run tests in the `backend/` directory to verify domain logic and atomicity:
```bash
cd backend
go test ./internal/test/...
```

---

For detailed technical notes on architecture, design decisions, and transaction management, refer to `notes.md`.
