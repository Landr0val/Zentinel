DO $migration$
DECLARE
    v_script_name VARCHAR := '000005_create_catalogues';
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
    CREATE TABLE IF NOT EXISTS catalogues (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        category VARCHAR(50) NOT NULL,
        code VARCHAR(50) NOT NULL,
        name VARCHAR(100) NOT NULL,
        description TEXT,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW(),
        UNIQUE (category, code)
    );

    CREATE INDEX IF NOT EXISTS idx_catalogues_category ON catalogues(category);

    -- Log Execution
    CALL sp_insert_replication_log(v_script_name, 'SUCCESS', 'Catalogues table created successfully');

EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Migration % failed: %', v_script_name, SQLERRM;
END $migration$;
