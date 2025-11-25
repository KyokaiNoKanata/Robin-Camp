package models

// Rating 评分模型
type Rating struct {
	MovieTitle string  `json:"movieTitle" db:"movie_title"`
	RaterID    string  `json:"raterId" db:"rater_id"`
	Rating     float64 `json:"rating" db:"rating"`
}

// RatingSubmit 提交评分请求
type RatingSubmit struct {
	Rating float64 `json:"rating" binding:"required,oneof=0.5 1.0 1.5 2.0 2.5 3.0 3.5 4.0 4.5 5.0"`
}

// RatingResult 评分结果响应
type RatingResult struct {
	MovieTitle string  `json:"movieTitle"`
	RaterID    string  `json:"raterId"`
	Rating     float64 `json:"rating"`
}

// RatingAggregate 评分聚合响应
type RatingAggregate struct {
	Average float64 `json:"average"`
	Count   int     `json:"count"`
}
