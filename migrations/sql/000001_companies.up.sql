BEGIN;

CREATE TABLE companies (
	company_id UUID DEFAULT gen_random_uuid ()
	,
	company_name VARCHAR NOT NULL
	,
	description VARCHAR NULL
	,
	amount_of_employees INT DEFAULT 0
	,
	is_registered BOOLEAN DEFAULT false
	,
    company_type VARCHAR NULL
    ,
	created_at timestamp without time zone
		NOT NULL
		DEFAULT CURRENT_TIMESTAMP
	,
	modified_at timestamp without time zone
	,
	deleted_at timestamp without time zone
	,
	PRIMARY KEY (company_id)
);

COMMIT;