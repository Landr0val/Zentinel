DO $migration$
DECLARE
    v_script_name VARCHAR := '000004_create_alerts';
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
    CREATE TABLE IF NOT EXISTS alerts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        transaction_id UUID NOT NULL REFERENCES transactions(id),
        client_id UUID NOT NULL REFERENCES clients(id),
        alert_type VARCHAR(50) NOT NULL,
        severity VARCHAR(20) NOT NULL,
        description TEXT,
        ai_explanation TEXT,
        status VARCHAR(20) DEFAULT 'pending',
        reviewed_by VARCHAR(100),
        reviewed_at TIMESTAMP,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );

    -- Log Execution
    CALL sp_insert_replication_log(v_script_name, 'SUCCESS', 'Alerts table created successfully');

EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Migration % failed: %', v_script_name, SQLERRM;
END $migration$;
