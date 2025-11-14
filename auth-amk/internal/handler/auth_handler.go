package handler

import (
	"net/http"

	"amk-backend/auth-amk/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler menampung endpoint utama autentikasi (login/refresh/logout).
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Identitas string `json:"identitas" binding:"required"`
	KataSandi string `json:"kata_sandi" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login menerima email/NRP dan kata sandi untuk menghasilkan token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "permintaan tidak valid"})
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Identitas, req.KataSandi, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case service.ErrKredensialSalah:
			status = http.StatusUnauthorized
		case service.ErrPenggunaNonAktif:
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Refresh membuat access token baru dari refresh token yang valid.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "permintaan tidak valid"})
		return
	}

	result, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrKredensialSalah {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Logout mencabut refresh token tertentu.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "permintaan tidak valid"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrKredensialSalah {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}
