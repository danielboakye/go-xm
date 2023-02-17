BEGIN;

ALTER TABLE companies
ADD CONSTRAINT companies_company_name_key UNIQUE (company_name);

COMMIT;