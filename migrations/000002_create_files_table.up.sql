CREATE TABLE files (
    id bigserial PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now(),
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    client_id INTEGER REFERENCES clients(id) ON DELETE CASCADE,
    category TEXT CHECK (category IN ('purchase_order', 'invoice', 'handing_over_report', 'pms_report'))
);
