package handler

import (
	"net/http"

	"amk-backend/auth-amk/internal/repository"
	"amk-backend/auth-amk/internal/service"

	"github.com/gin-gonic/gin"
)

// MeHandler mengembalikan profil pengguna yang sedang login.
type MeHandler struct {
	penggunaRepo *repository.PenggunaRepository
	hakService   *service.HakAksesService
}

func NewMeHandler(penggunaRepo *repository.PenggunaRepository, hakService *service.HakAksesService) *MeHandler {
	return &MeHandler{penggunaRepo: penggunaRepo, hakService: hakService}
}

func (h *MeHandler) Me(c *gin.Context) {
	userIDAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "token tidak valid"})
		return
	}
	userID := userIDAny.(uint)

	pengguna, err := h.penggunaRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	hak, err := h.hakService.ListByPenggunaID(c.Request.Context(), pengguna.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	peran := make([]string, 0, len(pengguna.Peran))
	for _, p := range pengguna.Peran {
		peran = append(peran, p.Nama)
	}

	c.JSON(http.StatusOK, gin.H{
		"pengguna":  pengguna,
		"peran":     peran,
		"hak_akses": hak,
	})
}
