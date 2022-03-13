BEGIN;
CREATE TABLE transactions(
    id serial PRIMARY KEY,
    user_id integer REFERENCES users ON DELETE CASCADE,
    order_id VARCHAR(255) NOT NULL,
    amount decimal NOT NULL,
    created_at timestamp
);
CREATE INDEX ON transactions (user_id);
COMMIT;