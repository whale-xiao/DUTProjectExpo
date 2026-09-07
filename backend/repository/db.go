package repository

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"showcase-backend/config"
	"showcase-backend/model"
)

// Open 建立 MySQL 连接（DSN 对齐本地 5.7 / 云端 8.0 均兼容）。
func Open(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

// Migrate 依据 model 自动建表（对应 03-数据库设计的新版 9 张表）。
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Project{},
		&model.Tag{},
		&model.ProjectTag{},
		&model.ProjectContentBlock{},
		&model.ProjectMember{},
		&model.ProjectImage{},
		&model.VisitLog{},
	)
}
