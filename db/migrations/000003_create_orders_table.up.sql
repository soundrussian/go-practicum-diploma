BEGIN;
CREATE TABLE orders(
    order_id varchar(255) PRIMARY KEY,
    user_id integer REFERENCES users ON DELETE CASCADE,
    accrual integer,
    status integer DEFAULT 0,
    uploaded_at timestamp
);
CREATE INDEX ON orders(user_id);
COMMIT;