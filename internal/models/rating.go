package models

// Rating 评分模型
type Rating struct {
	MovieTitle string  `json:"movieTitle" db:"movie_title"`
	RaterID    string  `json:"raterId" db:"rater_id"`
	Rating     float64 `json:"rating" db:"rating"`
}

// RatingSubmit 提交评分请求
type RatingSubmit struct {
	Score   float64 `json:"score" binding:"required,min=0,max=5"`
	Comment string  `json:"comment"`
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
