# Hexagonal Architecture Wallet API

This project implements a Wallet API using Hexagonal Architecture.

## Recent Changes

- Refactored `WalletService` interface and implementation to include `idempotencyKey`, `userID`, and `transactionID` for robust transaction handling.
- Updated `WalletHandler` to pass authenticated user ID and transaction ID to the service layer.
- Updated unit and integration tests to match the new service signature and ensure proper authorization and idempotency behavior.
