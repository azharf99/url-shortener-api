package config

import (
	"fmt"
	"os"
	"time"

	"github.com/azharf99/url-shortener-api/internal/domain"
	"github.com/azharf99/url-shortener-api/internal/repository"
	"github.com/azharf99/url-shortener-api/internal/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			zap.L().Warn("Error loading .env file")
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Set GORM logger to use Zap
	gormLogLevel := logger.Warn
	if os.Getenv("GIN_MODE") != "release" {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}

	db.AutoMigrate(&repository.UserModel{}, &repository.URLModel{})
	zap.L().Info("Database connected and migrated")

	seedAdmin(db)

	return db
}

func seedAdmin(db *gorm.DB) {
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" || adminEmail == "" || adminPassword == "" {
		zap.L().Warn("Skipping admin seeding: credentials not fully provided in env")
		return
	}

	var count int64
	db.Model(&repository.UserModel{}).Where("role = ?", domain.RoleAdmin).Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		zap.L().Error("Failed to hash admin password", zap.Error(err))
		return
	}

	admin := repository.UserModel{
		Username: adminUsername,
		Email:    adminEmail,
		Password: hashedPassword,
		Role:     domain.RoleAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&admin).Error; err != nil {
		zap.L().Error("Failed to seed admin", zap.Error(err))
	} else {
		zap.L().Info("Admin user seeded successfully", zap.String("username", adminUsername))
	}
}
