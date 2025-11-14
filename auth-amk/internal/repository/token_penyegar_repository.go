package repository

import (
	"context"
	"time"

	"amk-backend/auth-amk/internal/database"
	"amk-backend/auth-amk/internal/model"
)

// TokenPenyegarRepository mengelola data refresh token tersimpan di database.
type TokenPenyegarRepository struct{}

func NewTokenPenyegarRepository() *TokenPenyegarRepository { return &TokenPenyegarRepository{} }

func (r *TokenPenyegarRepository) Simpan(ctx context.Context, token *model.TokenPenyegar) error {
	return database.DB.WithContext(ctx).Create(token).Error
}

func (r *TokenPenyegarRepository) Perbarui(ctx context.Context, token *model.TokenPenyegar) error {
	return database.DB.WithContext(ctx).Save(token).Error
}

func (r *TokenPenyegarRepository) CariAktif(ctx context.Context, penggunaID uint, tokenHash string) (*model.TokenPenyegar, error) {
	var data model.TokenPenyegar
	err := database.DB.WithContext(ctx).
		Where("pengguna_id = ? AND token_hash = ?", penggunaID, tokenHash).
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *TokenPenyegarRepository) CariByHash(ctx context.Context, tokenHash string) (*model.TokenPenyegar, error) {
	var data model.TokenPenyegar
	err := database.DB.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *TokenPenyegarRepository) Cabut(ctx context.Context, tokenID uint) error {
	now := time.Now()
	return database.DB.WithContext(ctx).
		Model(&model.TokenPenyegar{}).
		Where("id = ?", tokenID).
		Updates(map[string]any{
			"dicabut_pada": now,
		}).Error
}

func (r *TokenPenyegarRepository) CabutSemuaPengguna(ctx context.Context, penggunaID uint) error {
	now := time.Now()
	return database.DB.WithContext(ctx).
		Model(&model.TokenPenyegar{}).
		Where("pengguna_id = ? AND dicabut_pada IS NULL", penggunaID).
		Updates(map[string]any{
			"dicabut_pada": now,
		}).Error
}
