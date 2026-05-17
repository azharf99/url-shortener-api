package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/azharf99/url-shortener-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type URLResponse struct {
	ID          uint      `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	ShortURL    string    `json:"short_url"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewURLResponse(u domain.URL) URLResponse {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8088/"
	}
	return URLResponse{
		ID:          u.ID,
		OriginalURL: u.OriginalURL,
		ShortCode:   u.ShortCode,
		ShortURL:    backendURL + u.ShortCode,
		UserID:      u.UserID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

type URLHandler struct {
	urlUsecase domain.URLUsecase
}

func NewURLHandler(u domain.URLUsecase) *URLHandler {
	return &URLHandler{u}
}

func (h *URLHandler) Shorten(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input struct {
		OriginalURL string `json:"original_url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.urlUsecase.Shorten(c.Request.Context(), userID, input.OriginalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create short URL"})
		return
	}

	c.JSON(http.StatusOK, NewURLResponse(*url))
}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	originalURL, err := h.urlUsecase.GetOriginalURL(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (h *URLHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	role := c.MustGet("role").(domain.Role)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	var input struct {
		OriginalURL string `json:"original_url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.urlUsecase.UpdateURL(c.Request.Context(), userID, role, uint(id), input.OriginalURL)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL updated successfully"})
}

func (h *URLHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	role := c.MustGet("role").(domain.Role)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	err = h.urlUsecase.DeleteURL(c.Request.Context(), userID, role, uint(id))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}

func (h *URLHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	role := c.MustGet("role").(domain.Role)

	urls, err := h.urlUsecase.ListURLs(c.Request.Context(), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not list URLs"})
		return
	}

	responses := make([]URLResponse, len(urls))
	for i, u := range urls {
		responses[i] = NewURLResponse(u)
	}

	c.JSON(http.StatusOK, gin.H{"urls": responses})
}
