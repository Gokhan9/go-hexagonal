CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    entity_id UUID NOT NULL, -- İşlem yapılan nesne ID'si (örn. wallet_id)
    entity_type VARCHAR(50) NOT NULL,  -- 'WALLET', 'TRANSACTION' vb.
    operation VARCHAR(50) NOT NULL, -- 'CREATE', 'UPDATE', 'TRANSFER'
    user_id VARCHAR(255), -- İşlemi yapan kullanıcı
    changes JSONB,  -- Değişen veriler (öncesi/sonrası veya metadata)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_id, entity_type);