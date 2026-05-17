package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type reCAPTCHAResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// Global HTTP client with timeout to prevent goroutine leaks
var captchaClient = &http.Client{
	Timeout: 10 * time.Second,
}

func CaptchaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		captchaToken := c.GetHeader("X-Recaptcha-Token")
		if captchaToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "reCAPTCHA token missing"})
			c.Abort()
			return
		}

		secret := os.Getenv("RECAPTCHA_SECRET")
		resp, err := captchaClient.PostForm("https://www.google.com/recaptcha/api/siteverify",
			url.Values{"secret": {secret}, "response": {captchaToken}})

		if err != nil {
			zap.L().Error("reCAPTCHA connection error", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify reCAPTCHA"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		var result reCAPTCHAResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			zap.L().Error("reCAPTCHA parse error", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reCAPTCHA response"})
			c.Abort()
			return
		}

		if !result.Success || result.Score < 0.5 {
			zap.L().Warn("reCAPTCHA blocked request", 
				zap.Bool("success", result.Success),
				zap.Float64("score", result.Score),
				zap.Strings("errors", result.ErrorCodes),
				zap.String("ip", c.ClientIP()),
			)
			c.JSON(http.StatusForbidden, gin.H{"error": "reCAPTCHA verification failed"})
			c.Abort()
			return
		}

		c.Next()
	}
}
