CREATE TABLE movie_stars (
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    star_id INTEGER NOT NULL REFERENCES stars(id),
    PRIMARY KEY (movie_id, star_id)
);

CREATE INDEX idx_movie_stars_movie_id ON movie_stars(movie_id);
CREATE INDEX idx_movie_stars_star_id ON movie_stars(star_id);

---- create above / drop below ----

DROP INDEX idx_movie_stars_star_id;
DROP INDEX idx_movie_stars_movie_id;
DROP TABLE movie_stars;
