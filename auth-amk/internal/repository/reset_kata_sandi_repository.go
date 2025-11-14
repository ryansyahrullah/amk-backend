package repository

import (
	"context"

	"amk-backend/auth-amk/internal/database"
	"amk-backend/auth-amk/internal/model"

	"gorm.io/gorm"
)

// ResetKataSandiRepository menangani penyimpanan kode OTP lupa sandi.
type ResetKataSandiRepository struct{}

func NewResetKataSandiRepository() *ResetKataSandiRepository { return &ResetKataSandiRepository{} }

func (r *ResetKataSandiRepository) Buat(ctx context.Context, data *model.ResetKataSandi) error {
	return database.DB.WithContext(ctx).Create(data).Error
}

func (r *ResetKataSandiRepository) CariAktif(ctx context.Context, penggunaID uint) (*model.ResetKataSandi, error) {
	var reset model.ResetKataSandi
	err := database.DB.WithContext(ctx).
		Where("pengguna_id = ? AND status = ?", penggunaID, "aktif").
		Order("id DESC").
		First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *ResetKataSandiRepository) Simpan(ctx context.Context, data *model.ResetKataSandi) error {
	return database.DB.WithContext(ctx).Save(data).Error
}

func (r *ResetKataSandiRepository) TambahPercobaan(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).
		Model(&model.ResetKataSandi{}).
		Where("id = ?", id).
		Update("jumlah_percobaan", gorm.Expr("jumlah_percobaan + 1")).Error
}

func (r *ResetKataSandiRepository) UbahStatus(ctx context.Context, id uint, status string) error {
	return database.DB.WithContext(ctx).
		Model(&model.ResetKataSandi{}).
		Where("id = ?", id).
		Update("status", status).Error
}
