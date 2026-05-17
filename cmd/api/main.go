package main

import (
	"github.com/azharf99/url-shortener-api/internal/config"
	"github.com/azharf99/url-shortener-api/internal/delivery/http/handler"
	"github.com/azharf99/url-shortener-api/internal/delivery/http/middleware"
	"github.com/azharf99/url-shortener-api/internal/repository"
	"github.com/azharf99/url-shortener-api/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()

	// Repositories
	urlRepo := repository.NewURLRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Usecases
	urlUsecase := usecase.NewURLUsecase(urlRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Handlers
	urlHandler := handler.NewURLHandler(urlUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	r.GET("/:shortCode", urlHandler.Redirect)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/shorten", urlHandler.Shorten)
		auth.GET("/urls", urlHandler.List)
		auth.PUT("/urls/:id", urlHandler.Update)
		auth.DELETE("/urls/:id", urlHandler.Delete)
	}

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.POST("/users", userHandler.Create)
		admin.GET("/users", userHandler.GetAll)
		admin.GET("/users/:id", userHandler.GetByID)
		admin.PUT("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id", userHandler.Delete)
	}

	r.Run(":8080")
}
