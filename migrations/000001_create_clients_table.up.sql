CREATE TABLE IF NOT EXISTS clients (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	company_name text NOT NULL,
	client_name text NOT NULL,
	email text NOT NULL,
	phone text NOT NULL
);
