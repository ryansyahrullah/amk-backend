package repository

import (
	"context"

	"amk-backend/auth-amk/internal/database"
	"amk-backend/auth-amk/internal/model"
)

// PeranRepository menangani query terkait peran.
type PeranRepository struct{}

func NewPeranRepository() *PeranRepository { return &PeranRepository{} }

// ListByPenggunaID mengembalikan seluruh peran milik pengguna.
func (r *PeranRepository) ListByPenggunaID(ctx context.Context, penggunaID uint) ([]model.Peran, error) {
	var peran []model.Peran
	err := database.DB.WithContext(ctx).
		Table("peran").
		Joins("JOIN pengguna_peran pp ON pp.peran_id = peran.id").
		Where("pp.pengguna_id = ?", penggunaID).
		Find(&peran).Error
	if err != nil {
		return nil, err
	}
	return peran, nil
}
