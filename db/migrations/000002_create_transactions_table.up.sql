BEGIN;
CREATE TABLE transactions(
    id serial PRIMARY KEY,
    user_id integer REFERENCES users ON DELETE CASCADE,
    amount integer NOT NULL,
    created_at timestamp
);
CREATE INDEX ON transactions (user_id);
COMMIT;