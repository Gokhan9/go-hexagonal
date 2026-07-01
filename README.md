# Go Hexagonal Wallet & Auth System

This project is a monorepo containing a robust, thread-safe Wallet and Authentication API implemented using **Hexagonal Architecture** (Ports & Adapters) in Go, accompanied by a modern **React/Vite** frontend.

## Project Structure

`
.
├── backend/    # Go Backend (Hexagonal Architecture)
└── frontend/   # React/TypeScript Frontend
`

## Backend (Go)

The backend follows the Hexagonal Architecture pattern to decouple core business logic from external technologies:

- **Core (Domain):** Entities (Wallet, Transaction, User), domain rules, and interfaces (Ports).
- **Service Layer:** Implements business logic (Wallet, User, JWT) using the defined Ports.
- **Adapters:**
  - **Driving (Primary):** handler package handles HTTP requests, including Auth and Middleware (Rate Limiting).
  - **Driven (Secondary):** epository package implements data storage (In-memory or PostgreSQL).

### Key Features
- **Concurrency Management:** Thread-safe operations using sync.Mutex and mechanisms to handle race conditions.
- **Idempotency:** Protects against duplicate requests using X-Idempotency-Key headers.
- **Authentication:** Secure user registration and login using JWT.
- **Rate Limiting:** Protects API endpoints from abuse.
- **Transaction Management:** Ensures data integrity for financial operations using database transactions (BeginTx, Commit, Rollback) and context-based propagation.
- **Persistence:** Supports In-memory and PostgreSQL storage.

### Getting Started (Backend)
Navigate to the ackend/ directory:
`ash
cd backend
go run cmd/main.go
`
The server will start on port 8080.

## Frontend (React/Vite)

A React application built with Vite and TypeScript for interacting with the backend API.

### Getting Started (Frontend)
Navigate to the rontend/ directory:
`ash
cd frontend
npm install
npm run dev
`

## Testing

### Backend
Run tests in the ackend/ directory:
`ash
cd backend
go test ./internal/test/...
`

---

For detailed technical notes on architecture and design decisions, refer to 
otes.md.
