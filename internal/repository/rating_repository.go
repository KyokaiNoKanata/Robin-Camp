package repository

import (
	"database/sql"
	"movie-rating-api/internal/models"
)

// RatingRepository 评分存储库接口
type RatingRepository interface {
	Upsert(rating *models.Rating) error
	GetByMovieAndRater(movieTitle, raterID string) (*models.Rating, error)
	GetAggregateByMovie(movieTitle string) (*models.RatingAggregate, error)
}

// ratingRepository 评分存储库实现
type ratingRepository struct {
	db *sql.DB
}

// NewRatingRepository 创建评分存储库实例
func NewRatingRepository(db *sql.DB) RatingRepository {
	return &ratingRepository{db: db}
}

// Upsert 更新或插入评分（Upsert操作）
func (r *ratingRepository) Upsert(rating *models.Rating) error {
	query := `
		INSERT INTO ratings (movie_title, rater_id, rating, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (movie_title, rater_id)
		DO UPDATE SET rating = $3, updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(query, rating.MovieTitle, rating.RaterID, rating.Rating)
	return err
}

// GetByMovieAndRater 根据电影标题和评分者ID获取评分
func (r *ratingRepository) GetByMovieAndRater(movieTitle, raterID string) (*models.Rating, error) {
	query := `
		SELECT movie_title, rater_id, rating
		FROM ratings
		WHERE movie_title = $1 AND rater_id = $2
	`

	var rating models.Rating
	err := r.db.QueryRow(query, movieTitle, raterID).Scan(
		&rating.MovieTitle, &rating.RaterID, &rating.Rating,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &rating, nil
}

// GetAggregateByMovie 获取电影的评分聚合信息（平均值和数量）
func (r *ratingRepository) GetAggregateByMovie(movieTitle string) (*models.RatingAggregate, error) {
	query := `
		SELECT ROUND(AVG(rating)::numeric, 1)::float, COUNT(*)
		FROM ratings
		WHERE movie_title = $1
	`

	var aggregate models.RatingAggregate
	err := r.db.QueryRow(query, movieTitle).Scan(
		&aggregate.Average, &aggregate.Count,
	)

	if err == sql.ErrNoRows {
		// 如果没有评分，返回平均值为0，数量为0
		return &models.RatingAggregate{Average: 0, Count: 0}, nil
	}
	if err != nil {
		return nil, err
	}

	return &aggregate, nil
}
