-- cüzdanlar
CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency VARCHAR(3) NOT NULL,
    version INT NOT NULL DEFAULT 1, -- Optimistic Locking
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- işlem geçmişi kayıtları
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallets(id),
    amount BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'DEPOSIT' veya 'WITHDRAW'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- idempotency (duplicate(yinelenen) istekleri engellemek için)
CREATE TABLE idempotency_records (
    idempotency_key VARCHAR(255) PRIMARY KEY,
    response_payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);