CREATE TABLE IF NOT EXISTS ratings (
    movie_title VARCHAR(255) NOT NULL,
    rater_id VARCHAR(255) NOT NULL,
    rating NUMERIC(2,1) NOT NULL CHECK (rating IN (0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (movie_title, rater_id),
    FOREIGN KEY (movie_title) REFERENCES movies(title) ON DELETE CASCADE
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_ratings_movie_title ON ratings(movie_title);
CREATE INDEX IF NOT EXISTS idx_ratings_rater_id ON ratings(rater_id);
