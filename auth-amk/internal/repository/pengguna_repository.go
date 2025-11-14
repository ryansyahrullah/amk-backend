package repository

import (
	"context"
	"errors"

	"amk-backend/auth-amk/internal/database"
	"amk-backend/auth-amk/internal/model"

	"gorm.io/gorm"
)

// ErrPenggunaTidakDitemukan dikembalikan ketika akun tidak ada di database.
var ErrPenggunaTidakDitemukan = errors.New("pengguna tidak ditemukan")

// PenggunaRepository menyediakan operasi terkait entitas pengguna.
type PenggunaRepository struct{}

func NewPenggunaRepository() *PenggunaRepository {
	return &PenggunaRepository{}
}

// FindByEmailOrNRP mencari pengguna berdasarkan email atau nrp.
func (r *PenggunaRepository) FindByEmailOrNRP(ctx context.Context, identitas string) (*model.Pengguna, error) {
	var pengguna model.Pengguna
	tx := database.DB.WithContext(ctx).
		Where("email = ? OR nrp = ?", identitas, identitas).
		Preload("Peran")
	if err := tx.First(&pengguna).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPenggunaTidakDitemukan
		}
		return nil, err
	}
	return &pengguna, nil
}

// FindByID mencari pengguna berdasarkan ID utama.
func (r *PenggunaRepository) FindByID(ctx context.Context, id uint) (*model.Pengguna, error) {
	var pengguna model.Pengguna
	if err := database.DB.WithContext(ctx).
		Preload("Peran").
		First(&pengguna, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPenggunaTidakDitemukan
		}
		return nil, err
	}
	return &pengguna, nil
}

// GetHakAksesByPenggunaID mengembalikan daftar kode hak akses yang dimiliki pengguna.
func (r *PenggunaRepository) GetHakAksesByPenggunaID(ctx context.Context, penggunaID uint) ([]string, error) {
	var hasil []string
	err := database.DB.WithContext(ctx).
		Table("hak_akses").
		Select("hak_akses.kode").
		Joins("JOIN peran_hak_akses pha ON pha.hak_akses_id = hak_akses.id").
		Joins("JOIN peran ON peran.id = pha.peran_id").
		Joins("JOIN pengguna_peran pp ON pp.peran_id = peran.id").
		Where("pp.pengguna_id = ?", penggunaID).
		Group("hak_akses.kode").
		Pluck("hak_akses.kode", &hasil).Error
	if err != nil {
		return nil, err
	}
	return hasil, nil
}

// UpdatePassword memperbarui hash kata sandi pengguna.
func (r *PenggunaRepository) UpdatePassword(ctx context.Context, penggunaID uint, hash string) error {
	return database.DB.WithContext(ctx).
		Model(&model.Pengguna{}).
		Where("id = ?", penggunaID).
		Update("kata_sandi_hash", hash).Error
}
