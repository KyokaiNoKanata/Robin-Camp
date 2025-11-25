package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"movie-rating-api/internal/models"
	"movie-rating-api/internal/service"
)

// MovieHandler 电影处理器
type MovieHandler struct {
	movieService  service.MovieService
	ratingService service.RatingService
}

// NewMovieHandler 创建电影处理器实例
func NewMovieHandler(movieService service.MovieService, ratingService service.RatingService) *MovieHandler {
	return &MovieHandler{
		movieService:  movieService,
		ratingService: ratingService,
	}
}

// CreateMovie 创建电影
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movieCreate models.MovieCreate

	// 绑定请求体
	if err := c.ShouldBindJSON(&movieCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 验证必填字段
	if movieCreate.Title == "" || movieCreate.ReleaseDate.IsZero() || len(movieCreate.Genre) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, release date and genre are required"})
		return
	}

	// 创建电影
	movie, err := h.movieService.CreateMovie(&movieCreate)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	// 设置Location头
	c.Header("Location", "/movies/"+movie.Title)
	c.JSON(http.StatusCreated, movie)
}

// ListMovies 列出电影
func (h *MovieHandler) ListMovies(c *gin.Context) {
	// 构建查询参数
	query := make(map[string]interface{})

	// 关键词搜索
	if q := c.Query("q"); q != "" {
		query["q"] = q
	}

	// 年份过滤
	if year := c.Query("year"); year != "" {
		query["year"] = year
	}

	// 类型过滤
	if genre := c.Query("genre"); genre != "" {
		query["genre"] = genre
	}

	// 分页参数
	limit := 10 // 默认值
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	cursor := c.Query("cursor")

	// 获取电影列表
	page, err := h.movieService.ListMovies(query, limit, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve movies"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// SubmitRating 提交评分
func (h *MovieHandler) SubmitRating(c *gin.Context) {
	// 获取路径参数
	movieTitle := c.Param("title")
	if movieTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie title is required"})
		return
	}

	var ratingSubmit models.RatingSubmit

	// 绑定请求体
	if err := c.ShouldBindJSON(&ratingSubmit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 设置电影标题
	ratingSubmit.MovieTitle = movieTitle

	// 验证必填字段
	if ratingSubmit.RaterID == "" || ratingSubmit.Rating < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rater ID and valid rating are required"})
		return
	}

	// 提交评分
	result, err := h.ratingService.SubmitRating(&ratingSubmit)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "rating must be") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit rating"})
		return
	}

	// 根据是否是更新操作返回不同的状态码
	statusCode := http.StatusCreated
	if result.Updated {
		statusCode = http.StatusOK
	}

	c.JSON(statusCode, result)
}

// GetMovieRatings 获取电影评分
func (h *MovieHandler) GetMovieRatings(c *gin.Context) {
	// 获取路径参数
	movieTitle := c.Param("title")
	if movieTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie title is required"})
		return
	}

	// 获取评分
	aggregate, err := h.ratingService.GetMovieRatings(movieTitle)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ratings"})
		return
	}

	c.JSON(http.StatusOK, aggregate)
}
