# Go项目Makefile

# 变量定义
APP_NAME := movie-rating-api
MAIN_FILE := ./cmd/api/main.go
BIN_DIR := ./bin
BIN_FILE := $(BIN_DIR)/$(APP_NAME)

# 默认目标
all: build

# 构建应用
build: 
	@echo "Building application..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_FILE) $(MAIN_FILE)
	@echo "Build completed: $(BIN_FILE)"

# 运行应用
run: build
	@echo "Running application..."
	@$(BIN_FILE)

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@echo "Clean completed"

# 测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 安装依赖
deps:
	@echo "Downloading dependencies..."
	@go mod tidy
	@go mod download

# 启动Docker容器
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

# 停止Docker容器
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

# 重启Docker容器
docker-restart:
	@echo "Restarting Docker containers..."
	@docker-compose down
	@docker-compose up -d

# 查看Docker日志
docker-logs:
	@echo "Showing Docker logs..."
	@docker-compose logs -f

# 查看API日志
docker-logs-api:
	@echo "Showing API logs..."
	@docker-compose logs -f api

# 查看数据库日志
docker-logs-db:
	@echo "Showing database logs..."
	@docker-compose logs -f db

# 检查容器状态
docker-status:
	@echo "Checking container status..."
	@docker-compose ps

# 帮助信息
help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Build and run the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make deps        - Download dependencies"
	@echo "  make docker-up   - Start Docker containers"
	@echo "  make docker-down - Stop Docker containers"
	@echo "  make docker-restart - Restart Docker containers"
	@echo "  make docker-logs - Show all Docker logs"
	@echo "  make docker-logs-api - Show API logs"
	@echo "  make docker-logs-db - Show database logs"
	@echo "  make docker-status - Check container status"

# 声明伪目标
.PHONY: all build run clean test fmt deps docker-up docker-down docker-restart docker-logs docker-logs-api docker-logs-db docker-status help
