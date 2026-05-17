package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter() gin.HandlerFunc {
	// Define the rate limit: 100 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	// Use in-memory store for now
	store := memory.NewStore()

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
