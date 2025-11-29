DO $migration$
DECLARE
    v_script_name VARCHAR := '000008_normalize_catalogues';
    v_already_executed BOOLEAN;
BEGIN
    -- 1. Check if already executed
    SELECT EXISTS(SELECT 1 FROM replication_logs WHERE script_name = v_script_name AND status = 'SUCCESS')
    INTO v_already_executed;

    IF v_already_executed THEN
        RAISE NOTICE 'Script % already executed. Skipping.', v_script_name;
        RETURN;
    END IF;

    -- 2. Create Master Data Types table
    EXECUTE $$
        CREATE TABLE IF NOT EXISTS master_data_types (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            code VARCHAR(50) NOT NULL UNIQUE,
            name VARCHAR(100) NOT NULL,
            description TEXT,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        )
    $$;

    -- 3. Seed Master Data Types
    EXECUTE $$
        INSERT INTO master_data_types (code, name, description) VALUES
        ('DOCUMENT_TYPE', 'Tipo de Documento', 'Identificación oficial'),
        ('RISK_PROFILE', 'Perfil de Riesgo', 'Clasificación de riesgo del cliente'),
        ('ACCOUNT_TYPE', 'Tipo de Cuenta', 'Categoría de producto bancario'),
        ('ACCOUNT_STATUS', 'Estado de Cuenta', 'Estado del ciclo de vida de la cuenta'),
        ('CURRENCY', 'Moneda', 'Divisas aceptadas'),
        ('OPERATION_TYPE', 'Tipo de Operación', 'Naturaleza de la transacción'),
        ('CHANNEL', 'Canal', 'Medio de realización de la operación'),
        ('TRANSACTION_STATUS', 'Estado de Transacción', 'Estado del procesamiento'),
        ('ALERT_TYPE', 'Tipo de Alerta', 'Categoría de la alerta de seguridad'),
        ('ALERT_SEVERITY', 'Severidad', 'Nivel de criticidad'),
        ('ALERT_STATUS', 'Estado de Alerta', 'Estado de gestión de la alerta')
        ON CONFLICT (code) DO NOTHING
    $$;

    -- 4. Modify Catalogues table
    EXECUTE $$ ALTER TABLE catalogues ADD COLUMN IF NOT EXISTS type_id UUID REFERENCES master_data_types(id) $$;

    -- Migrate existing data (if any) mapping category string to type_id
    EXECUTE $$
        UPDATE catalogues c
        SET type_id = mdt.id
        FROM master_data_types mdt
        WHERE c.category = mdt.code
    $$;

    -- If table was empty or had bad data, we clean up to enforce constraints
    EXECUTE $$ DELETE FROM catalogues WHERE type_id IS NULL $$;

    -- Drop old column and constraint
    EXECUTE $$ ALTER TABLE catalogues DROP CONSTRAINT IF EXISTS catalogues_category_code_key $$;
    EXECUTE $$ ALTER TABLE catalogues DROP COLUMN IF EXISTS category $$;
    EXECUTE $$ ALTER TABLE catalogues ALTER COLUMN type_id SET NOT NULL $$;

    -- Add new constraint
    EXECUTE $$ ALTER TABLE catalogues ADD CONSTRAINT catalogues_type_id_code_key UNIQUE (type_id, code) $$;

    -- 5. Re-Seed Catalogues (Upsert style to ensure data exists)

    -- DOCUMENT_TYPE
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('DNI', 'Documento Nacional de Identidad', 'Documento de identidad para ciudadanos'),
        ('PASSPORT', 'Pasaporte', 'Documento de viaje internacional'),
        ('CE', 'Carnet de Extranjería', 'Documento para residentes extranjeros'),
        ('RUC', 'Registro Único de Contribuyentes', 'Identificación tributaria')
    ) AS d(code, name, description)
    WHERE mdt.code = 'DOCUMENT_TYPE'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- RISK_PROFILE
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('low', 'Bajo Riesgo', 'Cliente con bajo riesgo de fraude'),
        ('standard', 'Riesgo Estándar', 'Perfil de riesgo promedio'),
        ('high', 'Alto Riesgo', 'Cliente con alto riesgo o PEP')
    ) AS d(code, name, description)
    WHERE mdt.code = 'RISK_PROFILE'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- ACCOUNT_TYPE
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('checking', 'Cuenta Corriente', 'Cuenta para transacciones diarias'),
        ('savings', 'Cuenta de Ahorros', 'Cuenta para ahorrar dinero'),
        ('credit_card', 'Tarjeta de Crédito', 'Línea de crédito rotativa')
    ) AS d(code, name, description)
    WHERE mdt.code = 'ACCOUNT_TYPE'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- ACCOUNT_STATUS
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('active', 'Activa', 'Cuenta operativa'),
        ('inactive', 'Inactiva', 'Sin movimiento reciente'),
        ('frozen', 'Congelada', 'Operaciones restringidas'),
        ('closed', 'Cerrada', 'Cuenta finalizada')
    ) AS d(code, name, description)
    WHERE mdt.code = 'ACCOUNT_STATUS'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- CURRENCY
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('USD', 'Dólar Estadounidense', 'Moneda de Estados Unidos'),
        ('EUR', 'Euro', 'Moneda de la Eurozona'),
        ('PEN', 'Sol Peruano', 'Moneda de Perú')
    ) AS d(code, name, description)
    WHERE mdt.code = 'CURRENCY'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- OPERATION_TYPE
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('purchase', 'Compra', 'Compra en comercio'),
        ('withdrawal', 'Retiro', 'Retiro de efectivo'),
        ('transfer', 'Transferencia', 'Transferencia entre cuentas'),
        ('deposit', 'Depósito', 'Ingreso de fondos')
    ) AS d(code, name, description)
    WHERE mdt.code = 'OPERATION_TYPE'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- CHANNEL
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('mobile', 'Banca Móvil', 'Aplicación móvil'),
        ('web', 'Banca por Internet', 'Portal web'),
        ('atm', 'Cajero Automático', 'ATM'),
        ('branch', 'Agencia', 'Atención presencial')
    ) AS d(code, name, description)
    WHERE mdt.code = 'CHANNEL'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- TRANSACTION_STATUS
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('pending', 'Pendiente', 'En proceso'),
        ('completed', 'Completada', 'Finalizada exitosamente'),
        ('failed', 'Fallida', 'Error en procesamiento'),
        ('reversed', 'Reversada', 'Operación anulada')
    ) AS d(code, name, description)
    WHERE mdt.code = 'TRANSACTION_STATUS'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- ALERT_TYPE
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('suspicious_activity', 'Actividad Sospechosa', 'Patrones inusuales'),
        ('high_amount', 'Monto Alto', 'Transacción supera umbrales'),
        ('location_mismatch', 'Ubicación Inusual', 'Ubicación diferente a la habitual'),
        ('velocity', 'Velocidad', 'Múltiples transacciones en poco tiempo')
    ) AS d(code, name, description)
    WHERE mdt.code = 'ALERT_TYPE'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- ALERT_SEVERITY
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('low', 'Baja', 'Riesgo menor'),
        ('medium', 'Media', 'Requiere revisión'),
        ('high', 'Alta', 'Acción inmediata requerida'),
        ('critical', 'Crítica', 'Bloqueo preventivo')
    ) AS d(code, name, description)
    WHERE mdt.code = 'ALERT_SEVERITY'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- ALERT_STATUS
    INSERT INTO catalogues (type_id, code, name, description)
    SELECT mdt.id, d.code, d.name, d.description
    FROM master_data_types mdt
    CROSS JOIN (VALUES
        ('pending', 'Pendiente', 'Sin revisar'),
        ('reviewing', 'En Revisión', 'Siendo analizada'),
        ('resolved', 'Resuelta', 'Caso cerrado'),
        ('false_positive', 'Falso Positivo', 'Alerta desestimada')
    ) AS d(code, name, description)
    WHERE mdt.code = 'ALERT_STATUS'
    ON CONFLICT (type_id, code) DO NOTHING;

    -- 6. Update Stored Procedures to match new schema
    EXECUTE $$ DROP PROCEDURE IF EXISTS sp_create_catalogue(VARCHAR, VARCHAR, VARCHAR, TEXT, BOOLEAN) $$;
    EXECUTE $$ DROP PROCEDURE IF EXISTS sp_update_catalogue(UUID, VARCHAR, VARCHAR, VARCHAR, TEXT, BOOLEAN) $$;

    EXECUTE $proc_def$
        CREATE OR REPLACE PROCEDURE sp_create_catalogue(
            IN p_type_code VARCHAR,
            IN p_code VARCHAR,
            IN p_name VARCHAR,
            IN p_description TEXT,
            IN p_is_active BOOLEAN
        )
        LANGUAGE plpgsql
        AS $inner$
        DECLARE
            v_catalogue_id UUID;
            v_type_id UUID;
        BEGIN
            SELECT id INTO v_type_id FROM master_data_types WHERE code = p_type_code;

            IF v_type_id IS NULL THEN
                RAISE EXCEPTION 'Master Data Type % not found', p_type_code;
            END IF;

            INSERT INTO catalogues (type_id, code, name, description, is_active)
            VALUES (v_type_id, p_code, p_name, p_description, p_is_active)
            RETURNING id INTO v_catalogue_id;

            CALL sp_insert_replication_log('sp_create_catalogue', 'SUCCESS', 'Catalogue created with ID: ' || v_catalogue_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_create_catalogue', 'ERROR', SQLERRM);
            RAISE;
        END;
        $inner$;
    $proc_def$;

    EXECUTE $proc_def$
        CREATE OR REPLACE PROCEDURE sp_update_catalogue(
            IN p_id UUID,
            IN p_type_code VARCHAR,
            IN p_code VARCHAR,
            IN p_name VARCHAR,
            IN p_description TEXT,
            IN p_is_active BOOLEAN
        )
        LANGUAGE plpgsql
        AS $inner$
        DECLARE
            v_type_id UUID;
        BEGIN
            IF p_type_code IS NOT NULL THEN
                SELECT id INTO v_type_id FROM master_data_types WHERE code = p_type_code;
                IF v_type_id IS NULL THEN
                    RAISE EXCEPTION 'Master Data Type % not found', p_type_code;
                END IF;
            END IF;

            UPDATE catalogues
            SET type_id = COALESCE(v_type_id, type_id),
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
        $inner$;
    $proc_def$;

    -- 7. Log execution
    CALL sp_insert_replication_log(v_script_name, 'SUCCESS', 'Migration executed successfully');

EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Migration % failed: %', v_script_name, SQLERRM;
END $migration$;
