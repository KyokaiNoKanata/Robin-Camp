package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"movie-rating-api/internal/models"
	"net/url"
	"time"
)

// BoxOfficeService 票房服务接口
type BoxOfficeService interface {
	GetBoxOfficeData(movieTitle string) (*models.BoxOffice, error)
}

// boxOfficeService 票房服务实现
type boxOfficeService struct {
	apiURL   string
	apiKey   string
	httpClient *http.Client
}

// NewBoxOfficeService 创建票房服务实例
func NewBoxOfficeService(apiURL, apiKey string) BoxOfficeService {
	if apiURL == "" || apiKey == "" {
		return nil
	}

	return &boxOfficeService{
		apiURL: apiURL,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBoxOfficeData 获取电影的票房数据
func (s *boxOfficeService) GetBoxOfficeData(movieTitle string) (*models.BoxOffice, error) {
	// 构建请求URL
	baseURL, err := url.Parse(s.apiURL)
	if err != nil {
		return nil, fmt.Errorf("invalid API URL: %v", err)
	}

	// 添加查询参数
	params := url.Values{}
	params.Add("title", movieTitle)
	params.Add("apikey", s.apiKey)
	baseURL.RawQuery = params.Encode()

	// 发送请求
	resp, err := s.httpClient.Get(baseURL.String())
	if err != nil {
		// 如果API调用失败，返回nil而不是错误，以便服务可以继续使用用户提供的数据
		return nil, nil
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	// 解析响应
	var boxOfficeResponse struct {
		Title       string `json:"title"`
		Distributor string `json:"distributor"`
		Budget      int64  `json:"budget"`
		MPARating   string `json:"mpa_rating"`
		Revenue     struct {
			Worldwide        int64 `json:"worldwide"`
			OpeningWeekendUSA int64 `json:"opening_weekend_usa"`
		} `json:"revenue"`
	}

	if err := json.Unmarshal(body, &boxOfficeResponse); err != nil {
		return nil, nil
	}

	// 构建BoxOffice数据
	return &models.BoxOffice{
		Revenue: models.Revenue{
			Worldwide:        boxOfficeResponse.Revenue.Worldwide,
			OpeningWeekendUSA: boxOfficeResponse.Revenue.OpeningWeekendUSA,
		},
		Distributor: boxOfficeResponse.Distributor,
		Budget:      boxOfficeResponse.Budget,
		MPARating:   boxOfficeResponse.MPARating,
		Currency:    "USD",
		Source:      "BoxOfficeAPI",
		LastUpdated: time.Now(),
	}, nil
}
