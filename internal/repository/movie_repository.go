package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"movie-rating-api/internal/models"
	"strings"
)

// MovieRepository 电影存储库接口
type MovieRepository interface {
	Create(movie *models.Movie) error
	GetByTitle(title string) (*models.Movie, error)
	List(query map[string]interface{}, limit int, cursor string) (*models.MoviePage, error)
	Update(movie *models.Movie) error
}

// movieRepository 电影存储库实现
type movieRepository struct {
	db *sql.DB
}

// NewMovieRepository 创建电影存储库实例
func NewMovieRepository(db *sql.DB) MovieRepository {
	return &movieRepository{db: db}
}

// Create 创建新电影
func (r *movieRepository) Create(movie *models.Movie) error {
	query := `
		INSERT INTO movies (id, title, release_date, genre, distributor, budget, mpa_rating, box_office)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var boxOfficeJSON []byte
	if movie.BoxOffice != nil {
		jsonData, err := json.Marshal(movie.BoxOffice)
		if err != nil {
			return err
		}
		boxOfficeJSON = jsonData
	}

	err := r.db.QueryRow(query, movie.ID, movie.Title, movie.ReleaseDate, movie.Genre,
		movie.Distributor, movie.Budget, movie.MPARating, boxOfficeJSON).Scan(&movie.ID)

	return err
}

// GetByTitle 根据标题获取电影
func (r *movieRepository) GetByTitle(title string) (*models.Movie, error) {
	query := `
		SELECT id, title, release_date, genre, distributor, budget, mpa_rating, box_office
		FROM movies
		WHERE title = $1
	`

	var movie models.Movie
	var boxOfficeJSON sql.NullString

	err := r.db.QueryRow(query, title).Scan(
		&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Genre,
		&movie.Distributor, &movie.Budget, &movie.MPARating, &boxOfficeJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 解析box_office JSON
	if boxOfficeJSON.Valid && boxOfficeJSON.String != "" {
		var boxOffice models.BoxOffice
		if err := json.Unmarshal([]byte(boxOfficeJSON.String), &boxOffice); err == nil {
			movie.BoxOffice = &boxOffice
		}
	}

	return &movie, nil
}

// List 列出电影，支持搜索和分页
func (r *movieRepository) List(query map[string]interface{}, limit int, cursor string) (*models.MoviePage, error) {
	if limit <= 0 {
		limit = 10
	}

	var conditions []string
	var args []interface{}
	argIndex := 1

	// 构建查询条件
	if q, ok := query["q"].(string); ok && q != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", argIndex))
		args = append(args, "%"+q+"%")
		argIndex++
	}

	if year, ok := query["year"].(int); ok && year > 0 {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM release_date) = $%d", argIndex))
		args = append(args, year)
		argIndex++
	}

	if genre, ok := query["genre"].(string); ok && genre != "" {
		conditions = append(conditions, fmt.Sprintf("genre ILIKE $%d", argIndex))
		args = append(args, genre)
		argIndex++
	}

	if distributor, ok := query["distributor"].(string); ok && distributor != "" {
		conditions = append(conditions, fmt.Sprintf("distributor ILIKE $%d", argIndex))
		args = append(args, distributor)
		argIndex++
	}

	if budget, ok := query["budget"].(int64); ok && budget > 0 {
		conditions = append(conditions, fmt.Sprintf("budget <= $%d", argIndex))
		args = append(args, budget)
		argIndex++
	}

	if mpaRating, ok := query["mpaRating"].(string); ok && mpaRating != "" {
		conditions = append(conditions, fmt.Sprintf("mpa_rating = $%d", argIndex))
		args = append(args, mpaRating)
		argIndex++
	}

	// 构建SQL查询
	sqlQuery := "SELECT id, title, release_date, genre, distributor, budget, mpa_rating, box_office FROM movies"
	if len(conditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	sqlQuery += " ORDER BY release_date DESC, title ASC"

	// 添加分页
	sqlQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, limit+1) // 获取多一行用于判断是否有下一页
	argIndex++

	// 执行查询
	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 解析结果
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		var boxOfficeJSON sql.NullString

		err := rows.Scan(
			&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Genre,
			&movie.Distributor, &movie.Budget, &movie.MPARating, &boxOfficeJSON,
		)
		if err != nil {
			return nil, err
		}

		// 解析box_office JSON
		if boxOfficeJSON.Valid && boxOfficeJSON.String != "" {
			var boxOffice models.BoxOffice
			if err := json.Unmarshal([]byte(boxOfficeJSON.String), &boxOffice); err == nil {
				movie.BoxOffice = &boxOffice
			}
		}

		movies = append(movies, movie)
	}

	// 构建分页响应
	result := &models.MoviePage{Items: movies}

	// 检查是否有下一页
	if len(movies) > limit {
		result.Items = movies[:limit]
		// 这里简单实现，实际应该返回基于ID或排序字段的游标
		nextCursor := "next"
		result.NextCursor = &nextCursor
	}

	return result, nil
}

// Update 更新电影信息
func (r *movieRepository) Update(movie *models.Movie) error {
	query := `
		UPDATE movies
		SET distributor = $1, budget = $2, mpa_rating = $3, box_office = $4
		WHERE title = $5
	`

	var boxOfficeJSON []byte
	if movie.BoxOffice != nil {
		jsonData, err := json.Marshal(movie.BoxOffice)
		if err != nil {
			return err
		}
		boxOfficeJSON = jsonData
	}

	_, err := r.db.Exec(query, movie.Distributor, movie.Budget, movie.MPARating, boxOfficeJSON, movie.Title)
	return err
}
