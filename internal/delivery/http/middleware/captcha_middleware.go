package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type reCAPTCHAResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
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
		resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
			url.Values{"secret": {secret}, "response": {captchaToken}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify reCAPTCHA"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		var result reCAPTCHAResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reCAPTCHA response"})
			c.Abort()
			return
		}

		if !result.Success || result.Score < 0.5 {
			c.JSON(http.StatusForbidden, gin.H{"error": "reCAPTCHA verification failed", "details": result.ErrorCodes})
			c.Abort()
			return
		}

		c.Next()
	}
}
