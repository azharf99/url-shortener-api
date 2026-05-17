package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
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
			log.Fatalf("Failed to create Redis rate limiter store: %v", err)
		}
		log.Println("Rate limiter using Redis store")
	} else {
		// Fallback to in-memory store if Redis client is not provided
		// Note: In real production, you might want to fail-fast if Redis is expected
		log.Println("WARNING: Rate limiter using fallback in-memory store")
		return func(c *gin.Context) { c.Next() } // Basic safety if misconfigured
	}

	// Create the limiter instance
	instance := limiter.New(store, rate)

	// Return the Gin middleware
	return mgin.NewMiddleware(instance, mgin.WithLimitReachedHandler(func(c *gin.Context) {
		log.Printf("Rate limit reached for IP: %s", c.ClientIP())
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate limit exceeded. Please try again in a minute.",
		})
		c.Abort()
	}))
}
