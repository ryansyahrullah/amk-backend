package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"amk-backend/auth-amk/internal/model"
	"amk-backend/auth-amk/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	// ErrKredensialSalah dikembalikan bila identitas atau kata sandi salah.
	ErrKredensialSalah = errors.New("identitas atau kata sandi salah")
	// ErrPenggunaNonAktif dikembalikan saat akun dinonaktifkan.
	ErrPenggunaNonAktif = errors.New("akun dalam status non_aktif")
)

// AuthService menangani alur login, refresh token, dan logout.
type AuthService struct {
	penggunaRepo *repository.PenggunaRepository
	peranRepo    *repository.PeranRepository
	tokenRepo    *repository.TokenPenyegarRepository
	tokenSvc     *TokenService
}

func NewAuthService(
	penggunaRepo *repository.PenggunaRepository,
	peranRepo *repository.PeranRepository,
	tokenRepo *repository.TokenPenyegarRepository,
	tokenSvc *TokenService,
) *AuthService {
	return &AuthService{
		penggunaRepo: penggunaRepo,
		peranRepo:    peranRepo,
		tokenRepo:    tokenRepo,
		tokenSvc:     tokenSvc,
	}
}

// LoginResult berisi token dan profil pengguna untuk response API.
type LoginResult struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	TipeToken    string          `json:"token_type"`
	Kadaluarsa   time.Time       `json:"refresh_expired_at"`
	Pengguna     *model.Pengguna `json:"pengguna"`
	Peran        []string        `json:"peran"`
	HakAkses     []string        `json:"hak_akses"`
}

// Login memvalidasi identitas dan kata sandi lalu mengembalikan token baru.
func (s *AuthService) Login(ctx context.Context, identitas, kataSandi, userAgent, ip string) (*LoginResult, error) {
	pengguna, err := s.penggunaRepo.FindByEmailOrNRP(ctx, identitas)
	if err != nil {
		if errors.Is(err, repository.ErrPenggunaTidakDitemukan) {
			return nil, ErrKredensialSalah
		}
		return nil, err
	}

	if !pengguna.Aktif() {
		return nil, ErrPenggunaNonAktif
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pengguna.KataSandiHash), []byte(kataSandi)); err != nil {
		return nil, ErrKredensialSalah
	}

	peranNama := make([]string, 0, len(pengguna.Peran))
	for _, p := range pengguna.Peran {
		peranNama = append(peranNama, p.Nama)
	}
	hakAkses, err := s.penggunaRepo.GetHakAksesByPenggunaID(ctx, pengguna.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(pengguna.ID, peranNama, hakAkses)
	if err != nil {
		return nil, err
	}

	plainRefresh, hashRefresh, expires, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenModel := &model.TokenPenyegar{
		PenggunaID:    pengguna.ID,
		TokenHash:     hashRefresh,
		UserAgent:     userAgent,
		AlamatIP:      ip,
		BerlakuSampai: expires,
	}
	if err := s.tokenRepo.Simpan(ctx, tokenModel); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: plainRefresh,
		TipeToken:    "Bearer",
		Kadaluarsa:   expires,
		Pengguna:     pengguna,
		Peran:        peranNama,
		HakAkses:     hakAkses,
	}, nil
}

// Refresh memperbarui access token menggunakan refresh token yang masih berlaku.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	hash := s.tokenSvc.HashRefreshToken(refreshToken)
	tokenModel, err := s.tokenRepo.CariByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKredensialSalah
		}
		return nil, err
	}

	if !tokenModel.Aktif(time.Now()) {
		return nil, fmt.Errorf("refresh token sudah tidak berlaku")
	}

	pengguna, err := s.penggunaRepo.FindByID(ctx, tokenModel.PenggunaID)
	if err != nil {
		if errors.Is(err, repository.ErrPenggunaTidakDitemukan) {
			return nil, ErrKredensialSalah
		}
		return nil, err
	}
	if !pengguna.Aktif() {
		return nil, ErrPenggunaNonAktif
	}

	peranList, err := s.peranRepo.ListByPenggunaID(ctx, pengguna.ID)
	if err != nil {
		return nil, err
	}
	peranNama := make([]string, 0, len(peranList))
	for _, p := range peranList {
		peranNama = append(peranNama, p.Nama)
	}
	hakAkses, err := s.penggunaRepo.GetHakAksesByPenggunaID(ctx, pengguna.ID)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.tokenSvc.GenerateAccessToken(pengguna.ID, peranNama, hakAkses)
	if err != nil {
		return nil, err
	}

	// Rotasi refresh token untuk keamanan
	plainRefresh, hashRefresh, expires, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	tokenModel.TokenHash = hashRefresh
	tokenModel.BerlakuSampai = expires
	tokenModel.DicabutPada = nil
	if err := s.tokenRepo.Perbarui(ctx, tokenModel); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: plainRefresh,
		TipeToken:    "Bearer",
		Kadaluarsa:   expires,
		Pengguna:     pengguna,
		Peran:        peranNama,
		HakAkses:     hakAkses,
	}, nil
}

// Logout mencabut refresh token tertentu.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hash := s.tokenSvc.HashRefreshToken(refreshToken)
	tokenModel, err := s.tokenRepo.CariByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrKredensialSalah
		}
		return err
	}
	return s.tokenRepo.Cabut(ctx, tokenModel.ID)
}

// CabutSemua sesi mencabut seluruh token refresh milik pengguna.
func (s *AuthService) CabutSemua(ctx context.Context, penggunaID uint) error {
	return s.tokenRepo.CabutSemuaPengguna(ctx, penggunaID)
}
