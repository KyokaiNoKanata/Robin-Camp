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

// Upsert 插入或更新评分
func (r *ratingRepository) Upsert(rating *models.Rating) error {
	// 由于我们已将数据库字段改为FLOAT类型，可以直接使用float64值
	query := `
		INSERT INTO ratings (movie_title, rater_id, rating, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (movie_title, rater_id)
		DO UPDATE SET rating = EXCLUDED.rating, updated_at = CURRENT_TIMESTAMP
	`

	// 直接执行，无需特殊处理
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

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rating, nil
}

// GetAggregateByMovie 获取电影的聚合评分
func (r *ratingRepository) GetAggregateByMovie(movieTitle string) (*models.RatingAggregate, error) {
	// 确保查询语句与FLOAT类型兼容
	query := `
		SELECT COALESCE(AVG(rating), 0) as average, COUNT(*) as count
		FROM ratings
		WHERE movie_title = $1
	`

	var aggregate models.RatingAggregate
	var avg float64
	var count int

	// 使用独立变量来扫描结果，确保类型兼容性
	err := r.db.QueryRow(query, movieTitle).Scan(&avg, &count)
	if err != nil {
		return nil, err
	}

	// 手动赋值给结构体
	aggregate.Average = avg
	aggregate.Count = count

	return &aggregate, nil
}
