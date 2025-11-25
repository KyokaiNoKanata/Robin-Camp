package main

import (
	"log"
	"net/http"
	"time"

	"movie-rating-api/internal/config"
	"movie-rating-api/internal/handlers"
	"movie-rating-api/internal/middleware"
	"movie-rating-api/internal/repository"
	"movie-rating-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 加载配置
	cfg := config.LoadConfig()
	log.Printf("Loaded AUTH_TOKEN: %s", cfg.AuthToken)

	// 强制使用本地PostgreSQL连接字符串
	localDBURL := "postgres://postgres:postgres@localhost:5432/movies?sslmode=disable"
	log.Printf("Using database URL: %s", localDBURL)

	// 初始化数据库
	db, err := repository.InitDB(localDBURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 确定迁移目录
	migrationDir := "internal/migrations"

	// 运行数据库迁移
	if err := repository.RunMigrations(localDBURL, migrationDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")

	// 初始化存储库
	movieRepo := repository.NewMovieRepository(db)
	ratingRepo := repository.NewRatingRepository(db)

	// 初始化服务
	boxOfficeService := service.NewBoxOfficeService(cfg.BoxOfficeURL, cfg.BoxOfficeAPIKey)
	movieService := service.NewMovieService(movieRepo, boxOfficeService)
	ratingService := service.NewRatingService(ratingRepo, movieRepo)

	// 初始化处理器
	movieHandler := handlers.NewMovieHandler(movieService, ratingService)
	healthHandler := handlers.NewHealthHandler()

	// 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// 设置Gin引擎
	router := gin.Default()

	// 注册中间件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 注册路由
	router.GET("/healthz", healthHandler.Check)

	// 需要认证的路由
	protected := router.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.POST("/movies", movieHandler.CreateMovie)
		protected.GET("/movies", movieHandler.ListMovies)
		protected.POST("/movies/:title/ratings", movieHandler.SubmitRating)
		protected.GET("/movies/:title/ratings", movieHandler.GetMovieRatings)
	}

	// 启动服务器 (使用端口9090)
	serverAddr := ":9090"
	log.Printf("Server starting on %s", serverAddr)

	// 使用自定义HTTP服务器配置
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
