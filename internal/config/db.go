package config

import (
	"fmt"
	"log"
	"os"

	"github.com/azharf99/url-shortener-api/internal/domain"
	"github.com/azharf99/url-shortener-api/internal/repository"
	"github.com/azharf99/url-shortener-api/internal/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file")
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&repository.UserModel{}, &repository.URLModel{})
	fmt.Println("Database connected and migrated")

	seedAdmin(db)

	return db
}

func seedAdmin(db *gorm.DB) {
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" || adminEmail == "" || adminPassword == "" {
		log.Println("Skipping admin seeding: credentials not fully provided in env")
		return
	}

	var count int64
	db.Model(&repository.UserModel{}).Where("role = ?", domain.RoleAdmin).Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	admin := repository.UserModel{
		Username: adminUsername,
		Email:    adminEmail,
		Password: hashedPassword,
		Role:     domain.RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to seed admin: %v", err)
	} else {
		fmt.Println("Admin user seeded successfully")
	}
}
