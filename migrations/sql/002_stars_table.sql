CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE stars (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    birth_date DATE NOT NULL,
    sex gender NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_stars_name ON stars (name);

---- create above / drop below ----

DROP INDEX idx_stars_name;
DROP TABLE stars;
DROP TYPE gender;