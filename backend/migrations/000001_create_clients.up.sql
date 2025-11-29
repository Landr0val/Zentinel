-- 1. Ensure replication infrastructure exists (Table and SP)
-- This is moved to migration 1 so all subsequent migrations can use it.
CREATE TABLE IF NOT EXISTS replication_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    script_name VARCHAR(255) NOT NULL,
    executed_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) NOT NULL,
    details TEXT
);

CREATE OR REPLACE PROCEDURE sp_insert_replication_log(
    IN p_script_name VARCHAR,
    IN p_status VARCHAR,
    IN p_details TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO replication_logs (script_name, status, details)
    VALUES (p_script_name, p_status, p_details);
END;
$$;

-- 2. Execute Migration Logic with Idempotency Check
DO $migration$
DECLARE
    v_script_name VARCHAR := '000001_create_clients';
    v_already_executed BOOLEAN;
BEGIN
    -- Check if already executed
    SELECT EXISTS(SELECT 1 FROM replication_logs WHERE script_name = v_script_name AND status = 'SUCCESS')
    INTO v_already_executed;

    IF v_already_executed THEN
        RAISE NOTICE 'Script % already executed. Skipping.', v_script_name;
        RETURN;
    END IF;

    -- Migration Logic: Create Clients Table
    CREATE TABLE IF NOT EXISTS clients (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        document_type VARCHAR(20) NOT NULL,
        document_number VARCHAR(50) NOT NULL UNIQUE,
        full_name VARCHAR(200) NOT NULL,
        email VARCHAR(150),
        phone VARCHAR(30),
        risk_profile VARCHAR(20) DEFAULT 'standard',
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );

    -- Log Execution
    CALL sp_insert_replication_log(v_script_name, 'SUCCESS', 'Clients table created successfully');

EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Migration % failed: %', v_script_name, SQLERRM;
END $migration$;
