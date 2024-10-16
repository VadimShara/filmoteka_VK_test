CREATE TYPE role AS ENUM ('user', 'admin');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(16) UNIQUE NOT NULL,
    pass_hash VARCHAR(64) NOT NULL,
    role role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_username ON users (username);

---- create above / drop below ----

DROP INDEX idx_users_username;
DROP TABLE users;
DROP TYPE role;