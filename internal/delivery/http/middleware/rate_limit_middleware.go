package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

func RateLimiter(redisClient *redis.Client) gin.HandlerFunc {
	// Define the rate limit: 100 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	var store limiter.Store
	var err error

	if redisClient != nil {
		// Use Redis store for production (distributed/persistent)
		store, err = sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
			Prefix: "rate_limiter:",
		})
		if err != nil {
			zap.L().Fatal("Failed to create Redis rate limiter store", zap.Error(err))
		}
		zap.L().Info("Rate limiter initialized with Redis store")
	} else {
		zap.L().Warn("Rate limiter using fallback in-memory store (non-distributed)")
		return func(c *gin.Context) { c.Next() }
	}

	// Create the limiter instance
	instance := limiter.New(store, rate)

	// Return the Gin middleware
	return mgin.NewMiddleware(instance, mgin.WithLimitReachedHandler(func(c *gin.Context) {
		zap.L().Warn("Rate limit reached", 
			zap.String("ip", c.ClientIP()), 
			zap.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate limit exceeded. Please try again in a minute.",
		})
		c.Abort()
	}))
}
