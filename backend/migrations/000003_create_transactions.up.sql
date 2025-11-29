CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    amount DECIMAL(18,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    operation_type VARCHAR(30) NOT NULL,
    channel VARCHAR(30) NOT NULL,
    merchant VARCHAR(200),
    country VARCHAR(3),
    city VARCHAR(100),
    risk_score INTEGER,
    is_flagged BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'completed',
    created_at TIMESTAMP DEFAULT NOW()
);
