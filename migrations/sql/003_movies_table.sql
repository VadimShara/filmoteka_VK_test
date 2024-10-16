CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    release_date DATE NOT NULL,
    rating INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_movies_title ON movies (title);
CREATE INDEX idx_movies_rating ON movies (rating);
CREATE INDEX idx_movies_release_date ON movies (release_date);

---- create above / drop below ----

DROP INDEX idx_movies_release_date;
DROP INDEX idx_movies_rating;
DROP INDEX idx_movies_title;
DROP TABLE movies;