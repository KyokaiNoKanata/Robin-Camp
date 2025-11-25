CREATE TABLE IF NOT EXISTS movies (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    genre VARCHAR(100) NOT NULL,
    distributor VARCHAR(255),
    budget BIGINT,
    mpa_rating VARCHAR(10),
    box_office JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(title)
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(title);
CREATE INDEX IF NOT EXISTS idx_movies_year ON EXTRACT(YEAR FROM release_date);
CREATE INDEX IF NOT EXISTS idx_movies_genre ON movies(genre);
CREATE INDEX IF NOT EXISTS idx_movies_distributor ON movies(distributor);
