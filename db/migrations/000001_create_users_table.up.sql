CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    login VARCHAR (255) UNIQUE NOT NULL,
    encrypted_password VARCHAR (255) NOT NULL
);