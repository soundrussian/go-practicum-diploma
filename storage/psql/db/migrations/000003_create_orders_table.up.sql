BEGIN;
CREATE TABLE orders(
    order_id varchar(255) PRIMARY KEY,
    user_id integer REFERENCES users ON DELETE CASCADE,
    accrual decimal NOT NULL DEFAULT 0,
    status integer DEFAULT 1,
    uploaded_at timestamp
);
CREATE INDEX ON orders(user_id);
COMMIT;