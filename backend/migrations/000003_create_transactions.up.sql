DO $migration$
DECLARE
    v_script_name VARCHAR := '000003_create_transactions';
    v_already_executed BOOLEAN;
BEGIN
    -- Check if already executed
    SELECT EXISTS(SELECT 1 FROM replication_logs WHERE script_name = v_script_name AND status = 'SUCCESS')
    INTO v_already_executed;

    IF v_already_executed THEN
        RAISE NOTICE 'Script % already executed. Skipping.', v_script_name;
        RETURN;
    END IF;

    -- Migration Logic
    CREATE TABLE IF NOT EXISTS transactions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        account_id UUID NOT NULL REFERENCES accounts(id),
        amount DECIMAL(18,2) NOT NULL,
        currency VARCHAR(3) DEFAULT 'USD',
        operation_type VARCHAR(50) NOT NULL,
        channel VARCHAR(50) NOT NULL,
        merchant VARCHAR(100),
        country VARCHAR(50),
        city VARCHAR(50),
        status VARCHAR(20) DEFAULT 'pending',
        risk_score INTEGER DEFAULT 0,
        is_flagged BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );

    -- Log Execution
    CALL sp_insert_replication_log(v_script_name, 'SUCCESS', 'Transactions table created successfully');

EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Migration % failed: %', v_script_name, SQLERRM;
END $migration$;
