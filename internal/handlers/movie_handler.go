package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"movie-rating-api/internal/models"
	"movie-rating-api/internal/service"

	"github.com/gin-gonic/gin"
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

// CreateMovie 创建新电影
func (h *MovieHandler) CreateMovie(c *gin.Context) {

	var movieCreate models.MovieCreate

	// 绑定请求体
	if err := c.ShouldBindJSON(&movieCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 验证必填字段
	if movieCreate.Title == "" || movieCreate.ReleaseDate == "" || movieCreate.Genre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, release date and genre are required"})
		return
	}
	
	// 验证ReleaseDate格式是否正确（YYYY-MM-DD）
	if _, err := time.Parse("2006-01-02", movieCreate.ReleaseDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Release date must be in YYYY-MM-DD format"})
		return
	}

	// 创建电影
	movie, err := h.movieService.CreateMovie(&movieCreate)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		// 记录详细错误信息
		fmt.Printf("Error creating movie: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create movie: %v", err)})
		return
	}

	// 设置Location头
	c.Header("Location", "/movies/"+movie.Title)
	c.JSON(http.StatusCreated, movie)
}

// ListMovies 获取电影列表
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

// SubmitRating 提交电影评分
func (h *MovieHandler) SubmitRating(c *gin.Context) {

	// 获取路径参数并进行URL解码
	movieTitle := c.Param("title")
	if movieTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie title is required"})
		return
	}
	// 解码URL中的'+'为空格
	movieTitle = strings.ReplaceAll(movieTitle, "+", " ")

	// 获取评分者ID（从查询参数或上下文）
	raterID := c.Query("raterId")
	if raterID == "" {
		// 尝试从请求头获取
		raterID = c.GetHeader("X-Rater-ID")
	}
	if raterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rater ID is required"})
		return
	}

	var ratingSubmit models.RatingSubmit

	// 绑定请求体
	if err := c.ShouldBindJSON(&ratingSubmit); err != nil {
		fmt.Printf("Invalid request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 验证评分
	if ratingSubmit.Score < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid rating is required"})
		return
	}

	// 提交评分
	result, err := h.ratingService.SubmitRating(movieTitle, raterID, &ratingSubmit)
	if err != nil {
		fmt.Printf("Error submitting rating: %v\n", err)
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

	// 返回评分结果
	c.JSON(http.StatusCreated, result)
}

// GetMovieRatings 获取电影评分
func (h *MovieHandler) GetMovieRatings(c *gin.Context) {

	// 获取路径参数并进行URL解码
	movieTitle := c.Param("title")
	if movieTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie title is required"})
		return
	}
	// 解码URL中的'+'为空格
	movieTitle = strings.ReplaceAll(movieTitle, "+", " ")

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
