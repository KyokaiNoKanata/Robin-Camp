package service

import (
	"fmt"
	"movie-rating-api/internal/models"
	"movie-rating-api/internal/repository"
	"strings"
	"time"
)

// MovieService 电影服务接口
type MovieService interface {
	CreateMovie(movieCreate *models.MovieCreate) (*models.Movie, error)
	GetMovieByTitle(title string) (*models.Movie, error)
	ListMovies(query map[string]interface{}, limit int, cursor string) (*models.MoviePage, error)
}

// movieService 电影服务实现
type movieService struct {
	movieRepo       repository.MovieRepository
	boxOfficeService BoxOfficeService
}

// NewMovieService 创建电影服务实例
func NewMovieService(movieRepo repository.MovieRepository, boxOfficeService BoxOfficeService) MovieService {
	return &movieService{
		movieRepo:       movieRepo,
		boxOfficeService: boxOfficeService,
	}
}

// CreateMovie 创建新电影
func (s *movieService) CreateMovie(movieCreate *models.MovieCreate) (*models.Movie, error) {
	// 检查电影是否已存在
	existingMovie, err := s.movieRepo.GetByTitle(movieCreate.Title)
	if err != nil {
		return nil, err
	}
	if existingMovie != nil {
		return nil, fmt.Errorf("movie with title '%s' already exists", movieCreate.Title)
	}

	// 创建电影实例
	movie := &models.Movie{
		ID:          generateMovieID(movieCreate.Title),
		Title:       movieCreate.Title,
		ReleaseDate: movieCreate.ReleaseDate,
		Genre:       movieCreate.Genre,
		Distributor: movieCreate.Distributor,
		Budget:      movieCreate.Budget,
		MPARating:   movieCreate.MPARating,
	}

	// 尝试从票房API获取数据
	if s.boxOfficeService != nil && movieCreate.Title != "" {
		boxOfficeData, err := s.boxOfficeService.GetBoxOfficeData(movieCreate.Title)
		if err == nil && boxOfficeData != nil {
			// 合并票房数据，但用户提供的值优先
			if movie.Distributor == nil && boxOfficeData.Distributor != "" {
				distributor := boxOfficeData.Distributor
				movie.Distributor = &distributor
			}

			if movie.Budget == nil && boxOfficeData.Budget != 0 {
				budget := boxOfficeData.Budget
				movie.Budget = &budget
			}

			if movie.MPARating == nil && boxOfficeData.MPARating != "" {
				mpaRating := boxOfficeData.MPARating
				movie.MPARating = &mpaRating
			}

			// 设置票房信息
			movie.BoxOffice = &models.BoxOffice{
				Revenue: models.Revenue{
					Worldwide:        boxOfficeData.Revenue.Worldwide,
					OpeningWeekendUSA: boxOfficeData.Revenue.OpeningWeekendUSA,
				},
				Currency:    "USD",
				Source:      "BoxOfficeAPI",
				LastUpdated: time.Now(),
			}
		}
	}

	// 保存到数据库
	if err := s.movieRepo.Create(movie); err != nil {
		return nil, err
	}

	return movie, nil
}

// GetMovieByTitle 根据标题获取电影
func (s *movieService) GetMovieByTitle(title string) (*models.Movie, error) {
	movie, err := s.movieRepo.GetByTitle(title)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, fmt.Errorf("movie not found")
	}
	return movie, nil
}

// ListMovies 列出电影
func (s *movieService) ListMovies(query map[string]interface{}, limit int, cursor string) (*models.MoviePage, error) {
	return s.movieRepo.List(query, limit, cursor)
}

// generateMovieID 生成电影ID
func generateMovieID(title string) string {
	// 简单的ID生成逻辑，可以根据需要改进
	cleanTitle := strings.ReplaceAll(strings.ToLower(title), " ", "_")
	return fmt.Sprintf("m_%s_%d", cleanTitle[:min(len(cleanTitle), 10)], time.Now().UnixNano()/1000000)
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
