package service

import (
	"fmt"
	"movie-rating-api/internal/models"
	"movie-rating-api/internal/repository"
)

// RatingService 评分服务接口
type RatingService interface {
	SubmitRating(submit *models.RatingSubmit) (*models.RatingResult, error)
	GetMovieRatings(movieTitle string) (*models.RatingAggregate, error)
}

// ratingService 评分服务实现
type ratingService struct {
	ratingRepo repository.RatingRepository
	movieRepo  repository.MovieRepository
}

// NewRatingService 创建评分服务实例
func NewRatingService(ratingRepo repository.RatingRepository, movieRepo repository.MovieRepository) RatingService {
	return &ratingService{
		ratingRepo: ratingRepo,
		movieRepo:  movieRepo,
	}
}

// SubmitRating 提交或更新评分
func (s *ratingService) SubmitRating(submit *models.RatingSubmit) (*models.RatingResult, error) {
	// 验证评分值
	if submit.Rating < 0 || submit.Rating > 10 {
		return nil, fmt.Errorf("rating must be between 0 and 10")
	}

	// 检查电影是否存在
	movie, err := s.movieRepo.GetByTitle(submit.MovieTitle)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, fmt.Errorf("movie not found")
	}

	// 创建评分实例
	rating := &models.Rating{
		MovieTitle: submit.MovieTitle,
		RaterID:    submit.RaterID,
		Rating:     submit.Rating,
	}

	// 检查是否是更新操作
	existingRating, err := s.ratingRepo.GetByMovieAndRater(submit.MovieTitle, submit.RaterID)
	if err != nil {
		return nil, err
	}

	// 保存评分
	if err := s.ratingRepo.Upsert(rating); err != nil {
		return nil, err
	}

	// 构建响应
	result := &models.RatingResult{
		MovieTitle: submit.MovieTitle,
		RaterID:    submit.RaterID,
		Rating:     submit.Rating,
		Updated:    true,
	}

	// 如果是新评分，Updated 为 false
	if existingRating == nil {
		result.Updated = false
	}

	return result, nil
}

// GetMovieRatings 获取电影的聚合评分
func (s *ratingService) GetMovieRatings(movieTitle string) (*models.RatingAggregate, error) {
	// 检查电影是否存在
	movie, err := s.movieRepo.GetByTitle(movieTitle)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, fmt.Errorf("movie not found")
	}

	// 获取聚合评分
	aggregate, err := s.ratingRepo.GetAggregateByMovie(movieTitle)
	if err != nil {
		return nil, err
	}

	return aggregate, nil
}
