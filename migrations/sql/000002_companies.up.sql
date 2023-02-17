BEGIN;

ALTER TABLE companies
DROP CONSTRAINT companies_company_name_key;

COMMIT;