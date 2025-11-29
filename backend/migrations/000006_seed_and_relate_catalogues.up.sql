DO $migration$
DECLARE
    v_script_name VARCHAR := '000006_seed_and_relate_catalogues';
    v_already_executed BOOLEAN;
    v_default_id UUID;
BEGIN
    -- Check if already executed
    SELECT EXISTS(SELECT 1 FROM replication_logs WHERE script_name = v_script_name AND status = 'SUCCESS')
    INTO v_already_executed;

    IF v_already_executed THEN
        RAISE NOTICE 'Script % already executed. Skipping.', v_script_name;
        RETURN;
    END IF;

    -- 1. Insert Seed Data
    INSERT INTO catalogues (category, code, name, description) VALUES
    -- DOCUMENT_TYPE
    ('DOCUMENT_TYPE', 'DNI', 'Documento Nacional de Identidad', 'Documento de identidad para ciudadanos'),
    ('DOCUMENT_TYPE', 'PASSPORT', 'Pasaporte', 'Documento de viaje internacional'),
    ('DOCUMENT_TYPE', 'CE', 'Carnet de Extranjería', 'Documento para residentes extranjeros'),
    ('DOCUMENT_TYPE', 'RUC', 'Registro Único de Contribuyentes', 'Identificación tributaria'),

    -- RISK_PROFILE
    ('RISK_PROFILE', 'low', 'Bajo Riesgo', 'Cliente con bajo riesgo de fraude'),
    ('RISK_PROFILE', 'standard', 'Riesgo Estándar', 'Perfil de riesgo promedio'),
    ('RISK_PROFILE', 'high', 'Alto Riesgo', 'Cliente con alto riesgo o PEP'),

    -- ACCOUNT_TYPE
    ('ACCOUNT_TYPE', 'checking', 'Cuenta Corriente', 'Cuenta para transacciones diarias'),
    ('ACCOUNT_TYPE', 'savings', 'Cuenta de Ahorros', 'Cuenta para ahorrar dinero'),
    ('ACCOUNT_TYPE', 'credit_card', 'Tarjeta de Crédito', 'Línea de crédito rotativa'),

    -- ACCOUNT_STATUS
    ('ACCOUNT_STATUS', 'active', 'Activa', 'Cuenta operativa'),
    ('ACCOUNT_STATUS', 'inactive', 'Inactiva', 'Sin movimiento reciente'),
    ('ACCOUNT_STATUS', 'frozen', 'Congelada', 'Operaciones restringidas'),
    ('ACCOUNT_STATUS', 'closed', 'Cerrada', 'Cuenta finalizada'),

    -- CURRENCY
    ('CURRENCY', 'USD', 'Dólar Estadounidense', 'Moneda de Estados Unidos'),
    ('CURRENCY', 'EUR', 'Euro', 'Moneda de la Eurozona'),
    ('CURRENCY', 'PEN', 'Sol Peruano', 'Moneda de Perú'),

    -- OPERATION_TYPE
    ('OPERATION_TYPE', 'purchase', 'Compra', 'Compra en comercio'),
    ('OPERATION_TYPE', 'withdrawal', 'Retiro', 'Retiro de efectivo'),
    ('OPERATION_TYPE', 'transfer', 'Transferencia', 'Transferencia entre cuentas'),
    ('OPERATION_TYPE', 'deposit', 'Depósito', 'Ingreso de fondos'),

    -- CHANNEL
    ('CHANNEL', 'mobile', 'Banca Móvil', 'Aplicación móvil'),
    ('CHANNEL', 'web', 'Banca por Internet', 'Portal web'),
    ('CHANNEL', 'atm', 'Cajero Automático', 'ATM'),
    ('CHANNEL', 'branch', 'Agencia', 'Atención presencial'),

    -- TRANSACTION_STATUS
    ('TRANSACTION_STATUS', 'pending', 'Pendiente', 'En proceso'),
    ('TRANSACTION_STATUS', 'completed', 'Completada', 'Finalizada exitosamente'),
    ('TRANSACTION_STATUS', 'failed', 'Fallida', 'Error en procesamiento'),
    ('TRANSACTION_STATUS', 'reversed', 'Reversada', 'Operación anulada'),

    -- ALERT_TYPE
    ('ALERT_TYPE', 'suspicious_activity', 'Actividad Sospechosa', 'Patrones inusuales'),
    ('ALERT_TYPE', 'high_amount', 'Monto Alto', 'Transacción supera umbrales'),
    ('ALERT_TYPE', 'location_mismatch', 'Ubicación Inusual', 'Ubicación diferente a la habitual'),
    ('ALERT_TYPE', 'velocity', 'Velocidad', 'Múltiples transacciones en poco tiempo'),

    -- ALERT_SEVERITY
    ('ALERT_SEVERITY', 'low', 'Baja', 'Riesgo menor'),
    ('ALERT_SEVERITY', 'medium', 'Media', 'Requiere revisión'),
    ('ALERT_SEVERITY', 'high', 'Alta', 'Acción inmediata requerida'),
    ('ALERT_SEVERITY', 'critical', 'Crítica', 'Bloqueo preventivo'),

    -- ALERT_STATUS
    ('ALERT_STATUS', 'pending', 'Pendiente', 'Sin revisar'),
    ('ALERT_STATUS', 'reviewing', 'En Revisión', 'Siendo analizada'),
    ('ALERT_STATUS', 'resolved', 'Resuelta', 'Caso cerrado'),
    ('ALERT_STATUS', 'false_positive', 'Falso Positivo', 'Alerta desestimada');

    -- 2. Relate Clients Table
    -- Document Type
    ALTER TABLE clients ADD COLUMN document_type_id UUID REFERENCES catalogues(id);
    UPDATE clients c SET document_type_id = cat.id FROM catalogues cat WHERE cat.code = c.document_type AND cat.category = 'DOCUMENT_TYPE';
    ALTER TABLE clients DROP COLUMN document_type;
    ALTER TABLE clients ALTER COLUMN document_type_id SET NOT NULL;

    -- Risk Profile
    ALTER TABLE clients ADD COLUMN risk_profile_id UUID REFERENCES catalogues(id);
    UPDATE clients c SET risk_profile_id = cat.id FROM catalogues cat WHERE cat.code = c.risk_profile AND cat.category = 'RISK_PROFILE';

    -- Set default for new rows (standard)
    SELECT id INTO v_default_id FROM catalogues WHERE category = 'RISK_PROFILE' AND code = 'standard';
    EXECUTE 'ALTER TABLE clients ALTER COLUMN risk_profile_id SET DEFAULT ' || quote_literal(v_default_id);

    ALTER TABLE clients DROP COLUMN risk_profile;


    -- 3. Relate Accounts Table
    -- Account Type
    ALTER TABLE accounts ADD COLUMN account_type_id UUID REFERENCES catalogues(id);
    UPDATE accounts a SET account_type_id = cat.id FROM catalogues cat WHERE cat.code = a.account_type AND cat.category = 'ACCOUNT_TYPE';
    ALTER TABLE accounts DROP COLUMN account_type;
    ALTER TABLE accounts ALTER COLUMN account_type_id SET NOT NULL;

    -- Account Status
    ALTER TABLE accounts ADD COLUMN status_id UUID REFERENCES catalogues(id);
    UPDATE accounts a SET status_id = cat.id FROM catalogues cat WHERE cat.code = a.status AND cat.category = 'ACCOUNT_STATUS';

    SELECT id INTO v_default_id FROM catalogues WHERE category = 'ACCOUNT_STATUS' AND code = 'active';
    EXECUTE 'ALTER TABLE accounts ALTER COLUMN status_id SET DEFAULT ' || quote_literal(v_default_id);

    ALTER TABLE accounts DROP COLUMN status;

    -- Currency
    ALTER TABLE accounts ADD COLUMN currency_id UUID REFERENCES catalogues(id);
    UPDATE accounts a SET currency_id = cat.id FROM catalogues cat WHERE cat.code = a.currency AND cat.category = 'CURRENCY';

    SELECT id INTO v_default_id FROM catalogues WHERE category = 'CURRENCY' AND code = 'USD';
    EXECUTE 'ALTER TABLE accounts ALTER COLUMN currency_id SET DEFAULT ' || quote_literal(v_default_id);

    ALTER TABLE accounts DROP COLUMN currency;


    -- 4. Relate Transactions Table
    -- Operation Type
    ALTER TABLE transactions ADD COLUMN operation_type_id UUID REFERENCES catalogues(id);
    UPDATE transactions t SET operation_type_id = cat.id FROM catalogues cat WHERE cat.code = t.operation_type AND cat.category = 'OPERATION_TYPE';
    ALTER TABLE transactions DROP COLUMN operation_type;
    ALTER TABLE transactions ALTER COLUMN operation_type_id SET NOT NULL;

    -- Channel
    ALTER TABLE transactions ADD COLUMN channel_id UUID REFERENCES catalogues(id);
    UPDATE transactions t SET channel_id = cat.id FROM catalogues cat WHERE cat.code = t.channel AND cat.category = 'CHANNEL';
    ALTER TABLE transactions DROP COLUMN channel;
    ALTER TABLE transactions ALTER COLUMN channel_id SET NOT NULL;

    -- Currency
    ALTER TABLE transactions ADD COLUMN currency_id UUID REFERENCES catalogues(id);
    UPDATE transactions t SET currency_id = cat.id FROM catalogues cat WHERE cat.code = t.currency AND cat.category = 'CURRENCY';

    SELECT id INTO v_default_id FROM catalogues WHERE category = 'CURRENCY' AND code = 'USD';
    EXECUTE 'ALTER TABLE transactions ALTER COLUMN currency_id SET DEFAULT ' || quote_literal(v_default_id);

    ALTER TABLE transactions DROP COLUMN currency;

    -- Status
    ALTER TABLE transactions ADD COLUMN status_id UUID REFERENCES catalogues(id);
    UPDATE transactions t SET status_id = cat.id FROM catalogues cat WHERE cat.code = t.status AND cat.category = 'TRANSACTION_STATUS';

    SELECT id INTO v_default_id FROM catalogues WHERE category = 'TRANSACTION_STATUS' AND code = 'completed';
    EXECUTE 'ALTER TABLE transactions ALTER COLUMN status_id SET DEFAULT ' || quote_literal(v_default_id);

    ALTER TABLE transactions DROP COLUMN status;


    -- 5. Relate Alerts Table
    -- Alert Type
    ALTER TABLE alerts ADD COLUMN alert_type_id UUID REFERENCES catalogues(id);
    UPDATE alerts a SET alert_type_id = cat.id FROM catalogues cat WHERE cat.code = a.alert_type AND cat.category = 'ALERT_TYPE';
    ALTER TABLE alerts DROP COLUMN alert_type;
    ALTER TABLE alerts ALTER COLUMN alert_type_id SET NOT NULL;

    -- Severity
    ALTER TABLE alerts ADD COLUMN severity_id UUID REFERENCES catalogues(id);
    UPDATE alerts a SET severity_id = cat.id FROM catalogues cat WHERE cat.code = a.severity AND cat.category = 'ALERT_SEVERITY';
    ALTER TABLE alerts DROP COLUMN severity;
    ALTER TABLE alerts ALTER COLUMN severity_id SET NOT NULL;

    -- Status
    ALTER TABLE alerts ADD COLUMN status_id UUID REFERENCES catalogues(id);
    UPDATE alerts a SET status_id = cat.id FROM catalogues cat WHERE cat.code = a.status AND cat.category = 'ALERT_STATUS';

    SELECT id INTO v_default_id FROM catalogues WHERE category = 'ALERT_STATUS' AND code = 'pending';
    EXECUTE 'ALTER TABLE alerts ALTER COLUMN status_id SET DEFAULT ' || quote_literal(v_default_id);

    ALTER TABLE alerts DROP COLUMN status;

    -- Log Execution
    CALL sp_insert_replication_log(v_script_name, 'SUCCESS', 'Catalogues seeded and tables related successfully');

EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Migration % failed: %', v_script_name, SQLERRM;
END $migration$;
