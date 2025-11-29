-- 1. Revert Alerts Table
-- Status
ALTER TABLE alerts ADD COLUMN status VARCHAR(20);
UPDATE alerts a SET status = cat.code FROM catalogues cat WHERE cat.id = a.status_id;
ALTER TABLE alerts ALTER COLUMN status SET DEFAULT 'pending';
ALTER TABLE alerts DROP COLUMN status_id;

-- Severity
ALTER TABLE alerts ADD COLUMN severity VARCHAR(20);
UPDATE alerts a SET severity = cat.code FROM catalogues cat WHERE cat.id = a.severity_id;
ALTER TABLE alerts ALTER COLUMN severity SET NOT NULL;
ALTER TABLE alerts DROP COLUMN severity_id;

-- Alert Type
ALTER TABLE alerts ADD COLUMN alert_type VARCHAR(50);
UPDATE alerts a SET alert_type = cat.code FROM catalogues cat WHERE cat.id = a.alert_type_id;
ALTER TABLE alerts ALTER COLUMN alert_type SET NOT NULL;
ALTER TABLE alerts DROP COLUMN alert_type_id;


-- 2. Revert Transactions Table
-- Status
ALTER TABLE transactions ADD COLUMN status VARCHAR(20);
UPDATE transactions t SET status = cat.code FROM catalogues cat WHERE cat.id = t.status_id;
ALTER TABLE transactions ALTER COLUMN status SET DEFAULT 'completed';
ALTER TABLE transactions DROP COLUMN status_id;

-- Currency
ALTER TABLE transactions ADD COLUMN currency VARCHAR(3);
UPDATE transactions t SET currency = cat.code FROM catalogues cat WHERE cat.id = t.currency_id;
ALTER TABLE transactions ALTER COLUMN currency SET DEFAULT 'USD';
ALTER TABLE transactions DROP COLUMN currency_id;

-- Channel
ALTER TABLE transactions ADD COLUMN channel VARCHAR(30);
UPDATE transactions t SET channel = cat.code FROM catalogues cat WHERE cat.id = t.channel_id;
ALTER TABLE transactions ALTER COLUMN channel SET NOT NULL;
ALTER TABLE transactions DROP COLUMN channel_id;

-- Operation Type
ALTER TABLE transactions ADD COLUMN operation_type VARCHAR(30);
UPDATE transactions t SET operation_type = cat.code FROM catalogues cat WHERE cat.id = t.operation_type_id;
ALTER TABLE transactions ALTER COLUMN operation_type SET NOT NULL;
ALTER TABLE transactions DROP COLUMN operation_type_id;


-- 3. Revert Accounts Table
-- Currency
ALTER TABLE accounts ADD COLUMN currency VARCHAR(3);
UPDATE accounts a SET currency = cat.code FROM catalogues cat WHERE cat.id = a.currency_id;
ALTER TABLE accounts ALTER COLUMN currency SET DEFAULT 'USD';
ALTER TABLE accounts DROP COLUMN currency_id;

-- Status
ALTER TABLE accounts ADD COLUMN status VARCHAR(20);
UPDATE accounts a SET status = cat.code FROM catalogues cat WHERE cat.id = a.status_id;
ALTER TABLE accounts ALTER COLUMN status SET DEFAULT 'active';
ALTER TABLE accounts DROP COLUMN status_id;

-- Account Type
ALTER TABLE accounts ADD COLUMN account_type VARCHAR(30);
UPDATE accounts a SET account_type = cat.code FROM catalogues cat WHERE cat.id = a.account_type_id;
ALTER TABLE accounts ALTER COLUMN account_type SET NOT NULL;
ALTER TABLE accounts DROP COLUMN account_type_id;


-- 4. Revert Clients Table
-- Risk Profile
ALTER TABLE clients ADD COLUMN risk_profile VARCHAR(20);
UPDATE clients c SET risk_profile = cat.code FROM catalogues cat WHERE cat.id = c.risk_profile_id;
ALTER TABLE clients ALTER COLUMN risk_profile SET DEFAULT 'standard';
ALTER TABLE clients DROP COLUMN risk_profile_id;

-- Document Type
ALTER TABLE clients ADD COLUMN document_type VARCHAR(20);
UPDATE clients c SET document_type = cat.code FROM catalogues cat WHERE cat.id = c.document_type_id;
ALTER TABLE clients ALTER COLUMN document_type SET NOT NULL;
ALTER TABLE clients DROP COLUMN document_type_id;


-- 5. Clean Catalogues
DELETE FROM catalogues;
