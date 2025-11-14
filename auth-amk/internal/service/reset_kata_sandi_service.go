package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"amk-backend/auth-amk/internal/model"
	"amk-backend/auth-amk/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ResetKataSandiService mengelola alur lupa kata sandi berbasis OTP.
type ResetKataSandiService struct {
	penggunaRepo *repository.PenggunaRepository
	resetRepo    *repository.ResetKataSandiRepository
	tokenRepo    *repository.TokenPenyegarRepository
}

func NewResetKataSandiService(
	penggunaRepo *repository.PenggunaRepository,
	resetRepo *repository.ResetKataSandiRepository,
	tokenRepo *repository.TokenPenyegarRepository,
) *ResetKataSandiService {
	return &ResetKataSandiService{
		penggunaRepo: penggunaRepo,
		resetRepo:    resetRepo,
		tokenRepo:    tokenRepo,
	}
}

// MintaOTP membuat kode OTP baru dan menyimpannya di database.
func (s *ResetKataSandiService) MintaOTP(ctx context.Context, email string) (string, error) {
	pengguna, err := s.penggunaRepo.FindByEmailOrNRP(ctx, email)
	if err != nil {
		return "", err
	}

	otp, err := generateOTP()
	if err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	data := &model.ResetKataSandi{
		PenggunaID:      pengguna.ID,
		KodeOTPHash:     string(hash),
		BerlakuSampai:   time.Now().Add(10 * time.Minute),
		JumlahPercobaan: 0,
		Status:          "aktif",
	}
	if err := s.resetRepo.Buat(ctx, data); err != nil {
		return "", err
	}

	return otp, nil
}

// KonfirmasiOTP memvalidasi kode OTP dan mengubah kata sandi pengguna.
func (s *ResetKataSandiService) KonfirmasiOTP(ctx context.Context, email, otpBaru, kataSandiBaru string) error {
	pengguna, err := s.penggunaRepo.FindByEmailOrNRP(ctx, email)
	if err != nil {
		return err
	}

	record, err := s.resetRepo.CariAktif(ctx, pengguna.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("kode OTP tidak ditemukan")
		}
		return err
	}

	now := time.Now()
	if now.After(record.BerlakuSampai) {
		_ = s.resetRepo.UbahStatus(ctx, record.ID, "kedaluwarsa")
		return fmt.Errorf("kode OTP kedaluwarsa")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(record.KodeOTPHash), []byte(otpBaru)); err != nil {
		_ = s.resetRepo.TambahPercobaan(ctx, record.ID)
		if record.JumlahPercobaan+1 >= 5 {
			_ = s.resetRepo.UbahStatus(ctx, record.ID, "kedaluwarsa")
			return fmt.Errorf("kode OTP salah dan diblokir")
		}
		return fmt.Errorf("kode OTP tidak cocok")
	}

	hashBaru, err := bcrypt.GenerateFromPassword([]byte(kataSandiBaru), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.penggunaRepo.UpdatePassword(ctx, pengguna.ID, string(hashBaru)); err != nil {
		return err
	}

	if err := s.resetRepo.UbahStatus(ctx, record.ID, "dipakai"); err != nil {
		return err
	}

	// Cabut semua refresh token supaya sesi lain ikut logout.
	if err := s.tokenRepo.CabutSemuaPengguna(ctx, pengguna.ID); err != nil {
		return err
	}

	return nil
}

func generateOTP() (string, error) {
	const digits = "0123456789"
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b), nil
}
