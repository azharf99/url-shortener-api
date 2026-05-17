package main

import (
	"os"
	"strings"

	"github.com/azharf99/url-shortener-api/internal/config"
	"github.com/azharf99/url-shortener-api/internal/delivery/http/handler"
	"github.com/azharf99/url-shortener-api/internal/delivery/http/middleware"
	"github.com/azharf99/url-shortener-api/internal/repository"
	"github.com/azharf99/url-shortener-api/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize Logger
	config.InitLogger()
	defer zap.L().Sync()

	db := config.ConnectDB()
	redisClient := config.ConnectRedis()

	// Repositories
	urlRepo := repository.NewURLRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Usecases
	urlUsecase := usecase.NewURLUsecase(urlRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Handlers
	urlHandler := handler.NewURLHandler(urlUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New() // Use gin.New() to avoid default logger
	r.Use(middleware.ZapLogger(), gin.Recovery())

	// Rate Limiter (using Redis if available)
	r.Use(middleware.RateLimiter(redisClient))

	// CORS
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	origins := []string{"*"}
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Recaptcha-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	r.GET("/:shortCode", urlHandler.Redirect)
	
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/google-login", userHandler.GoogleLogin)

		// Protected routes
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware(), middleware.CaptchaMiddleware())
		{
			auth.POST("/shorten", urlHandler.Shorten)
			auth.GET("/urls", urlHandler.List)
			auth.PUT("/urls/:id", urlHandler.Update)
			auth.DELETE("/urls/:id", urlHandler.Delete)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly(), middleware.CaptchaMiddleware())
		{
			admin.POST("/users", userHandler.Create)
			admin.GET("/users", userHandler.List)
			admin.GET("/users/:id", userHandler.GetByID)
			admin.PUT("/users/:id", userHandler.Update)
			admin.DELETE("/users/:id", userHandler.Delete)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	zap.L().Info("Server starting", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		zap.L().Fatal("Server failed to start", zap.Error(err))
	}
}
