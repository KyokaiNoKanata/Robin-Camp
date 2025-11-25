package models

import (
	"encoding/json"
	"time"
)

// Movie 电影模型
type Movie struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title" binding:"required"`
	ReleaseDate string     `json:"releaseDate" db:"release_date" binding:"required"`
	Genre       string     `json:"genre" db:"genre" binding:"required"`
	Distributor *string    `json:"distributor,omitempty" db:"distributor"`
	Budget      *int64     `json:"budget,omitempty" db:"budget"`
	MPARating   *string    `json:"mpaRating,omitempty" db:"mpa_rating"`
	BoxOffice   *BoxOffice `json:"boxOffice,omitempty" db:"box_office"`
}

// BoxOffice 票房信息
type BoxOffice struct {
	Revenue     Revenue  `json:"revenue"`
	Currency    string   `json:"currency"`
	Source      string   `json:"source"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// Revenue 收入信息
type Revenue struct {
	Worldwide        int64  `json:"worldwide"`
	OpeningWeekendUSA *int64 `json:"openingWeekendUSA,omitempty"`
}

// BoxOfficeJSON 用于数据库存储的票房JSON格式
type BoxOfficeJSON struct {
	Revenue     map[string]interface{} `json:"revenue"`
	Currency    string                 `json:"currency"`
	Source      string                 `json:"source"`
	LastUpdated time.Time              `json:"lastUpdated"`
}

// Value 实现driver.Valuer接口，用于存储到数据库
func (b *BoxOffice) Value() (interface{}, error) {
	if b == nil {
		return nil, nil
	}

	// 转换为JSON格式
	boJSON := BoxOfficeJSON{
		Currency:    b.Currency,
		Source:      b.Source,
		LastUpdated: b.LastUpdated,
		Revenue: map[string]interface{}{
			"worldwide": b.Revenue.Worldwide,
		},
	}

	if b.Revenue.OpeningWeekendUSA != nil {
		boJSON.Revenue["openingWeekendUSA"] = *b.Revenue.OpeningWeekendUSA
	}

	return json.Marshal(boJSON)
}

// Scan 实现sql.Scanner接口，用于从数据库读取
func (b *BoxOffice) Scan(value interface{}) error {
	if value == nil {
		b = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	var boJSON BoxOfficeJSON
	if err := json.Unmarshal(bytes, &boJSON); err != nil {
		return err
	}

	b.Currency = boJSON.Currency
	b.Source = boJSON.Source
	b.LastUpdated = boJSON.LastUpdated

	// 解析收入信息
	worldwide, _ := boJSON.Revenue["worldwide"].(float64)
	b.Revenue.Worldwide = int64(worldwide)

	if openingWeekendUSA, ok := boJSON.Revenue["openingWeekendUSA"]; ok {
		val := int64(openingWeekendUSA.(float64))
		b.Revenue.OpeningWeekendUSA = &val
	}

	return nil
}

// MovieCreate 创建电影请求
type MovieCreate struct {
	Title       string  `json:"title" binding:"required"`
	ReleaseDate string  `json:"releaseDate" binding:"required"`
	Genre       string  `json:"genre" binding:"required"`
	Distributor *string `json:"distributor,omitempty"`
	Budget      *int64  `json:"budget,omitempty"`
	MPARating   *string `json:"mpaRating,omitempty"`
}

// MoviePage 电影分页响应
type MoviePage struct {
	Items      []Movie `json:"items"`
	NextCursor *string `json:"nextCursor,omitempty"`
}
