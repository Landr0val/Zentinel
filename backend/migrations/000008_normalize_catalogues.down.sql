-- 1. Revert Stored Procedures
CREATE OR REPLACE PROCEDURE sp_create_catalogue(
    IN p_category VARCHAR,
    IN p_code VARCHAR,
    IN p_name VARCHAR,
    IN p_description TEXT,
    IN p_is_active BOOLEAN
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_catalogue_id UUID;
BEGIN
    INSERT INTO catalogues (category, code, name, description, is_active)
    VALUES (p_category, p_code, p_name, p_description, p_is_active)
    RETURNING id INTO v_catalogue_id;

    CALL sp_insert_replication_log('sp_create_catalogue', 'SUCCESS', 'Catalogue created with ID: ' || v_catalogue_id);
EXCEPTION WHEN OTHERS THEN
    CALL sp_insert_replication_log('sp_create_catalogue', 'ERROR', SQLERRM);
    RAISE;
END;
$$;

CREATE OR REPLACE PROCEDURE sp_update_catalogue(
    IN p_id UUID,
    IN p_category VARCHAR,
    IN p_code VARCHAR,
    IN p_name VARCHAR,
    IN p_description TEXT,
    IN p_is_active BOOLEAN
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE catalogues
    SET category = p_category,
        code = p_code,
        name = p_name,
        description = p_description,
        is_active = p_is_active,
        updated_at = NOW()
    WHERE id = p_id;

    CALL sp_insert_replication_log('sp_update_catalogue', 'SUCCESS', 'Catalogue updated with ID: ' || p_id);
EXCEPTION WHEN OTHERS THEN
    CALL sp_insert_replication_log('sp_update_catalogue', 'ERROR', SQLERRM);
    RAISE;
END;
$$;

-- 2. Revert Catalogues Table Structure
-- Add category column
ALTER TABLE catalogues ADD COLUMN category VARCHAR(50);

-- Migrate data back from type_id to category string
UPDATE catalogues c
SET category = mdt.code
FROM master_data_types mdt
WHERE c.type_id = mdt.id;

-- Enforce NOT NULL on category
ALTER TABLE catalogues ALTER COLUMN category SET NOT NULL;

-- Drop new unique constraint
ALTER TABLE catalogues DROP CONSTRAINT IF EXISTS catalogues_type_id_code_key;

-- Drop type_id column (and implicitly the FK constraint)
ALTER TABLE catalogues DROP COLUMN type_id;

-- Restore old unique constraint
ALTER TABLE catalogues ADD CONSTRAINT catalogues_category_code_key UNIQUE (category, code);

-- 3. Drop Master Data Types Table
DROP TABLE IF EXISTS master_data_types;
