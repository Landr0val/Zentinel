DO $migration$
DECLARE
    v_script_name VARCHAR := '000007_create_replication_log_and_sps';
    v_already_executed BOOLEAN;
BEGIN
    -- Check if already executed
    SELECT EXISTS(SELECT 1 FROM replication_logs WHERE script_name = v_script_name AND status = 'SUCCESS')
    INTO v_already_executed;

    IF v_already_executed THEN
        RAISE NOTICE 'Script % already executed. Skipping.', v_script_name;
        RETURN;
    END IF;

    -- Migration Logic: Create Stored Procedures

    -- sp_create_catalogue
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_create_catalogue(
            IN p_category VARCHAR,
            IN p_code VARCHAR,
            IN p_name VARCHAR,
            IN p_description TEXT,
            IN p_is_active BOOLEAN
        )
        LANGUAGE plpgsql
        AS $body$
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
        $body$;
    $proc$;

    -- sp_create_client
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_create_client(
            IN p_document_type_id UUID,
            IN p_document_number VARCHAR,
            IN p_full_name VARCHAR,
            IN p_email VARCHAR,
            IN p_phone VARCHAR,
            IN p_risk_profile_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        DECLARE
            v_client_id UUID;
        BEGIN
            INSERT INTO clients (document_type_id, document_number, full_name, email, phone, risk_profile_id)
            VALUES (p_document_type_id, p_document_number, p_full_name, p_email, p_phone, p_risk_profile_id)
            RETURNING id INTO v_client_id;

            CALL sp_insert_replication_log('sp_create_client', 'SUCCESS', 'Client created with ID: ' || v_client_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_create_client', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_create_account
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_create_account(
            IN p_client_id UUID,
            IN p_account_number VARCHAR,
            IN p_account_type_id UUID,
            IN p_currency_id UUID,
            IN p_balance DECIMAL,
            IN p_status_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        DECLARE
            v_account_id UUID;
        BEGIN
            INSERT INTO accounts (client_id, account_number, account_type_id, currency_id, balance, status_id)
            VALUES (p_client_id, p_account_number, p_account_type_id, p_currency_id, p_balance, p_status_id)
            RETURNING id INTO v_account_id;

            CALL sp_insert_replication_log('sp_create_account', 'SUCCESS', 'Account created with ID: ' || v_account_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_create_account', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_create_transaction
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_create_transaction(
            IN p_account_id UUID,
            IN p_amount DECIMAL,
            IN p_currency_id UUID,
            IN p_operation_type_id UUID,
            IN p_channel_id UUID,
            IN p_merchant VARCHAR,
            IN p_country VARCHAR,
            IN p_city VARCHAR,
            IN p_status_id UUID,
            IN p_risk_score INTEGER,
            IN p_is_flagged BOOLEAN
        )
        LANGUAGE plpgsql
        AS $body$
        DECLARE
            v_transaction_id UUID;
        BEGIN
            INSERT INTO transactions (account_id, amount, currency_id, operation_type_id, channel_id, merchant, country, city, status_id, risk_score, is_flagged)
            VALUES (p_account_id, p_amount, p_currency_id, p_operation_type_id, p_channel_id, p_merchant, p_country, p_city, p_status_id, p_risk_score, p_is_flagged)
            RETURNING id INTO v_transaction_id;

            CALL sp_insert_replication_log('sp_create_transaction', 'SUCCESS', 'Transaction created with ID: ' || v_transaction_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_create_transaction', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_create_alert
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_create_alert(
            IN p_transaction_id UUID,
            IN p_client_id UUID,
            IN p_alert_type_id UUID,
            IN p_severity_id UUID,
            IN p_description TEXT,
            IN p_ai_explanation TEXT,
            IN p_status_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        DECLARE
            v_alert_id UUID;
        BEGIN
            INSERT INTO alerts (transaction_id, client_id, alert_type_id, severity_id, description, ai_explanation, status_id)
            VALUES (p_transaction_id, p_client_id, p_alert_type_id, p_severity_id, p_description, p_ai_explanation, p_status_id)
            RETURNING id INTO v_alert_id;

            CALL sp_insert_replication_log('sp_create_alert', 'SUCCESS', 'Alert created with ID: ' || v_alert_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_create_alert', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_update_catalogue
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_update_catalogue(
            IN p_id UUID,
            IN p_category VARCHAR,
            IN p_code VARCHAR,
            IN p_name VARCHAR,
            IN p_description TEXT,
            IN p_is_active BOOLEAN
        )
        LANGUAGE plpgsql
        AS $body$
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
        $body$;
    $proc$;

    -- sp_delete_catalogue
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_delete_catalogue(
            IN p_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            DELETE FROM catalogues WHERE id = p_id;

            CALL sp_insert_replication_log('sp_delete_catalogue', 'SUCCESS', 'Catalogue deleted with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_delete_catalogue', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_update_client
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_update_client(
            IN p_id UUID,
            IN p_document_type_id UUID,
            IN p_document_number VARCHAR,
            IN p_full_name VARCHAR,
            IN p_email VARCHAR,
            IN p_phone VARCHAR,
            IN p_risk_profile_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            UPDATE clients
            SET document_type_id = p_document_type_id,
                document_number = p_document_number,
                full_name = p_full_name,
                email = p_email,
                phone = p_phone,
                risk_profile_id = p_risk_profile_id,
                updated_at = NOW()
            WHERE id = p_id;

            CALL sp_insert_replication_log('sp_update_client', 'SUCCESS', 'Client updated with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_update_client', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_delete_client
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_delete_client(
            IN p_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            DELETE FROM clients WHERE id = p_id;

            CALL sp_insert_replication_log('sp_delete_client', 'SUCCESS', 'Client deleted with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_delete_client', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_update_account
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_update_account(
            IN p_id UUID,
            IN p_client_id UUID,
            IN p_account_number VARCHAR,
            IN p_account_type_id UUID,
            IN p_currency_id UUID,
            IN p_balance DECIMAL,
            IN p_status_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            UPDATE accounts
            SET client_id = p_client_id,
                account_number = p_account_number,
                account_type_id = p_account_type_id,
                currency_id = p_currency_id,
                balance = p_balance,
                status_id = p_status_id,
                updated_at = NOW()
            WHERE id = p_id;

            CALL sp_insert_replication_log('sp_update_account', 'SUCCESS', 'Account updated with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_update_account', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_delete_account
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_delete_account(
            IN p_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            DELETE FROM accounts WHERE id = p_id;

            CALL sp_insert_replication_log('sp_delete_account', 'SUCCESS', 'Account deleted with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_delete_account', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_update_transaction
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_update_transaction(
            IN p_id UUID,
            IN p_account_id UUID,
            IN p_amount DECIMAL,
            IN p_currency_id UUID,
            IN p_operation_type_id UUID,
            IN p_channel_id UUID,
            IN p_merchant VARCHAR,
            IN p_country VARCHAR,
            IN p_city VARCHAR,
            IN p_status_id UUID,
            IN p_risk_score INTEGER,
            IN p_is_flagged BOOLEAN
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            UPDATE transactions
            SET account_id = p_account_id,
                amount = p_amount,
                currency_id = p_currency_id,
                operation_type_id = p_operation_type_id,
                channel_id = p_channel_id,
                merchant = p_merchant,
                country = p_country,
                city = p_city,
                status_id = p_status_id,
                risk_score = p_risk_score,
                is_flagged = p_is_flagged
            WHERE id = p_id;

            CALL sp_insert_replication_log('sp_update_transaction', 'SUCCESS', 'Transaction updated with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_update_transaction', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_delete_transaction
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_delete_transaction(
            IN p_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            DELETE FROM transactions WHERE id = p_id;

            CALL sp_insert_replication_log('sp_delete_transaction', 'SUCCESS', 'Transaction deleted with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_delete_transaction', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_update_alert
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_update_alert(
            IN p_id UUID,
            IN p_transaction_id UUID,
            IN p_client_id UUID,
            IN p_alert_type_id UUID,
            IN p_severity_id UUID,
            IN p_description TEXT,
            IN p_ai_explanation TEXT,
            IN p_status_id UUID,
            IN p_reviewed_by VARCHAR,
            IN p_reviewed_at TIMESTAMP
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            UPDATE alerts
            SET transaction_id = p_transaction_id,
                client_id = p_client_id,
                alert_type_id = p_alert_type_id,
                severity_id = p_severity_id,
                description = p_description,
                ai_explanation = p_ai_explanation,
                status_id = p_status_id,
                reviewed_by = p_reviewed_by,
                reviewed_at = p_reviewed_at
            WHERE id = p_id;

            CALL sp_insert_replication_log('sp_update_alert', 'SUCCESS', 'Alert updated with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_update_alert', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- sp_delete_alert
    EXECUTE $proc$
        CREATE OR REPLACE PROCEDURE sp_delete_alert(
            IN p_id UUID
        )
        LANGUAGE plpgsql
        AS $body$
        BEGIN
            DELETE FROM alerts WHERE id = p_id;

            CALL sp_insert_replication_log('sp_delete_alert', 'SUCCESS', 'Alert deleted with ID: ' || p_id);
        EXCEPTION WHEN OTHERS THEN
            CALL sp_insert_replication_log('sp_delete_alert', 'ERROR', SQLERRM);
            RAISE;
        END;
        $body$;
    $proc$;

    -- Log Execution
    CALL sp_insert_replication_log(v_script_name, 'SUCCESS', 'Stored Procedures created successfully');

EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Migration % failed: %', v_script_name, SQLERRM;
END $migration$;
