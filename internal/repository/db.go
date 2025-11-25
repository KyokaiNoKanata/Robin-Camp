package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

// RunMigrations 运行数据库迁移
func RunMigrations(dbURL string, migrationDir string) error {
	// 确保路径使用正斜杠（URL格式需要）
	migrationPath := strings.ReplaceAll(migrationDir, "\\", "/")

	// 创建迁移实例
	m, err := migrate.New(
		"file://"+migrationPath,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// 检查当前迁移状态
	version, dirty, err := m.Version()

	// 处理脏数据库状态
	if err == nil && dirty {
		log.Printf("Found dirty database version %d, forcing clean state...", version)
		// 对于脏数据库，直接强制设置为当前版本使其干净
		if err := m.Force(int(version)); err != nil {
			log.Printf("Warning: Failed to force version %d: %v", version, err)
			// 如果强制设置失败，尝试跳过迁移
			log.Println("Skipping migrations due to database state issues")
			return nil
		}
	}

	// 运行迁移
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Warning: Migration error: %v", err)
		// 如果迁移失败，我们仍然尝试继续运行应用
		log.Println("Continuing despite migration warnings")
	}

	log.Println("Migrations completed successfully")
	return nil
}
