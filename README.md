# Go Hexagonal Wallet & Auth API

This project demonstrates a robust, thread-safe Wallet and Authentication API implemented using **Hexagonal Architecture** (Ports & Adapters) in Go. It focuses on domain-driven design, transactional integrity, concurrency management, idempotency, authentication (JWT), and rate limiting.

## Architecture

This project follows the Hexagonal Architecture pattern to decouple core business logic from external technologies:

- **Core (Domain):** Entities (`Wallet`, `Transaction`, `User`), domain rules, and interfaces (Ports).
- **Service Layer:** Implements business logic (Wallet, User, JWT) using the defined Ports.
- **Adapters:**
  - **Driving (Primary):** `handler` package handles HTTP requests, including Auth and Middleware (Rate Limiting).
  - **Driven (Secondary):** `repository` package implements data storage (In-memory or PostgreSQL).

## Key Features

- **Concurrency Management:** Thread-safe operations using `sync.Mutex` and mechanisms to handle race conditions.
- **Idempotency:** Protects against duplicate requests using `X-Idempotency-Key` headers.
- **Authentication:** Secure user registration and login using JWT.
- **Rate Limiting:** Protects API endpoints from abuse.
- **Transaction Management:** Ensures data integrity for financial operations using database transactions (`BeginTx`, `Commit`, `Rollback`) and context-based propagation.
- **Domain-Driven Design:** Strong domain rules and guard clauses to ensure financial integrity.
- **Persistence:** Supports In-memory and PostgreSQL storage.
- **Comprehensive Testing:** Unit and integration testing.

## Getting Started

To run the application:

```bash
go run cmd/main.go
```

The server will start on port `8080`.

## Testing with Postman

### API Endpoints

**Authentication:**
- **Register:** `POST /register`
- **Login:** `POST /login`

**Wallets:**
- **Create Wallet:** `POST /wallets`
- **Get Wallet:** `GET /wallets/{id}`
- **Deposit:** `POST /wallets/{id}/deposit`
- **Withdraw:** `POST /wallets/{id}/withdraw`
- **Get Transactions:** `GET /wallets/{id}/transactions`

### Postman Instructions

1.  **Base URL:** `http://localhost:8080`
2.  **Headers:**
    - For transactional operations (Deposit/Withdraw), you **must** provide the `X-Idempotency-Key` header to prevent duplicate processing.
      - Key: `X-Idempotency-Key`
      - Value: A unique UUID (e.g., `550e8400-e29b-41d4-a716-446655440000`)
    - For protected endpoints, provide the `Authorization` header with the JWT token received from `/login`:
      - Key: `Authorization`
      - Value: `Bearer <your_token>`
3.  **JSON Body:** Example for Deposit/Withdraw:
    ```json
    {
      "amount": 100.0,
      "user_id": "user123"
    }
    ```

For detailed technical notes on architecture and design decisions, refer to `notes.md`.
