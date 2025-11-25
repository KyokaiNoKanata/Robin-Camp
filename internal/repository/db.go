package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// InitDB 初始化数据库连接
func InitDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

// RunMigrations 执行数据库迁移
func RunMigrations(dbURL string) error {
	// 打开数据库连接
	db, err := InitDB(dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// 创建迁移源
	migrationDir := filepath.Join("internal", "migrations")

	// 检查迁移目录是否存在
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		// 如果在Docker容器中运行，尝试使用绝对路径
		migrationDir = "/app/internal/migrations"
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			return fmt.Errorf("migrations directory not found")
		}
	}

	// 创建迁移源
	source, err := file.New(migrationDir, "")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	// 创建PostgreSQL驱动
	driver, err := postgres.WithInstance(db, &postgres.Config{SchemaName: "public"})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// 创建迁移实例
	m, err := migrate.NewWithInstance("file", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// 执行迁移
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
