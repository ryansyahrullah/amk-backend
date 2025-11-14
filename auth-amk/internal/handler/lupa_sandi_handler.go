package handler

import (
	"net/http"

	"amk-backend/auth-amk/internal/service"

	"github.com/gin-gonic/gin"
)

// LupaSandiHandler menangani permintaan lupa sandi (OTP dan reset kata sandi).
type LupaSandiHandler struct {
	resetService *service.ResetKataSandiService
}

func NewLupaSandiHandler(resetService *service.ResetKataSandiService) *LupaSandiHandler {
	return &LupaSandiHandler{resetService: resetService}
}

type mintaOTPRequest struct {
	Email string `json:"email" binding:"required"`
}

type resetSandiRequest struct {
	Email         string `json:"email" binding:"required"`
	KodeOTP       string `json:"kode_otp" binding:"required"`
	KataSandiBaru string `json:"kata_sandi_baru" binding:"required"`
}

func (h *LupaSandiHandler) MintaOTP(c *gin.Context) {
	var req mintaOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "permintaan tidak valid"})
		return
	}

	otp, err := h.resetService.MintaOTP(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// OTP dikembalikan di response untuk kebutuhan pengujian (di produksi akan dikirim via email).
	c.JSON(http.StatusOK, gin.H{
		"message": "kode OTP dikirim ke email",
		"otp_dev": otp,
	})
}

func (h *LupaSandiHandler) ResetKataSandi(c *gin.Context) {
	var req resetSandiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "permintaan tidak valid"})
		return
	}

	if err := h.resetService.KonfirmasiOTP(c.Request.Context(), req.Email, req.KodeOTP, req.KataSandiBaru); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kata sandi berhasil diperbarui"})
}
