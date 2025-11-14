package service

import (
	"context"

	"amk-backend/auth-amk/internal/repository"
)

// HakAksesService menyediakan helper untuk mengambil daftar hak akses.
type HakAksesService struct {
	penggunaRepo *repository.PenggunaRepository
}

func NewHakAksesService(penggunaRepo *repository.PenggunaRepository) *HakAksesService {
	return &HakAksesService{penggunaRepo: penggunaRepo}
}

func (s *HakAksesService) ListByPenggunaID(ctx context.Context, penggunaID uint) ([]string, error) {
	return s.penggunaRepo.GetHakAksesByPenggunaID(ctx, penggunaID)
}
