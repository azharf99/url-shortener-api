package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/azharf99/url-shortener-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	userRepo domain.UserRepository
}

func NewSubscriptionHandler(userRepo domain.UserRepository) *SubscriptionHandler {
	return &SubscriptionHandler{userRepo}
}

func (h *SubscriptionHandler) Checkout(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	// Fetch user details
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	pgURL := os.Getenv("PAYMENT_GATEWAY_URL")
	pgAPIKey := os.Getenv("PAYMENT_GATEWAY_API_KEY")
	if pgURL == "" || pgAPIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment gateway configuration is missing on server"})
		return
	}

	orderID := fmt.Sprintf("sub-%d-%d", userID, time.Now().Unix())

	// Call payment-gateway-api CreateCharge
	chargePayload := map[string]interface{}{
		"order_id":     orderID,
		"gross_amount": 100000,
		"customer_details": map[string]interface{}{
			"first_name": user.Username,
			"email":      user.Email,
		},
		"credit_card": map[string]interface{}{
			"save_card": true,
		},
	}
	payloadBytes, _ := json.Marshal(chargePayload)

	req, err := http.NewRequest("POST", pgURL+"/api/v1/charge", bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to construct request to payment gateway"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", pgAPIKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach payment gateway: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var pgErr struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&pgErr)
		c.JSON(resp.StatusCode, gin.H{"error": "Payment gateway error: " + pgErr.Error})
		return
	}

	var chargeResp struct {
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
		OrderID     string `json:"order_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&chargeResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment gateway response"})
		return
	}

	c.JSON(http.StatusOK, chargeResp)
}

type ForwardedNotification struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	SavedTokenID      string `json:"saved_token_id"`
	SubscriptionID    string `json:"subscription_id"`
}

func (h *SubscriptionHandler) HandlePaymentWebhook(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	expectedKey := os.Getenv("PAYMENT_GATEWAY_API_KEY")
	if apiKey == "" || apiKey != expectedKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid X-API-Key"})
		return
	}

	var notif ForwardedNotification
	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification payload"})
		return
	}

	// OrderID format: sub-<userID>-<timestamp>
	parts := strings.Split(notif.OrderID, "-")
	if len(parts) < 2 || parts[0] != "sub" {
		c.JSON(http.StatusOK, gin.H{"message": "Not a subscription order ID"})
		return
	}

	userIDVal, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in order ID"})
		return
	}
	userID := uint(userIDVal)

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User not found, ignoring webhook"})
		return
	}

	// Process payment status
	status := notif.TransactionStatus
	if status == "settlement" || status == "capture" {
		user.IsPremium = true
		// Grant 30 days subscription from now, or extend if currently active
		if user.SubscriptionEnd.After(time.Now()) {
			user.SubscriptionEnd = user.SubscriptionEnd.AddDate(0, 1, 0)
		} else {
			user.SubscriptionEnd = time.Now().AddDate(0, 1, 0)
		}

		if err := h.userRepo.Update(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
			return
		}

		// If this is the initial checkout and we got a saved token, create a recurring subscription schedule
		if notif.SavedTokenID != "" && notif.SubscriptionID == "" {
			h.createRecurringSubscription(user.ID, notif.SavedTokenID)
		}
	} else if status == "deny" || status == "expire" || status == "cancel" {
		// If recurring payment failed, set IsPremium = false
		if notif.SubscriptionID != "" {
			user.IsPremium = false
			_ = h.userRepo.Update(ctx, user)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *SubscriptionHandler) createRecurringSubscription(userID uint, savedTokenID string) {
	pgURL := os.Getenv("PAYMENT_GATEWAY_URL")
	pgAPIKey := os.Getenv("PAYMENT_GATEWAY_API_KEY")
	if pgURL == "" || pgAPIKey == "" {
		return
	}

	payload := map[string]interface{}{
		"user_id":        userID,
		"saved_token_id": savedTokenID,
		"gross_amount":   100000,
		"name":           "Premium Subscription Monthly",
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", pgURL+"/api/v1/subscriptions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", pgAPIKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}
