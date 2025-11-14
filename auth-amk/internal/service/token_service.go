package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"amk-backend/auth-amk/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// AccessClaims adalah payload JWT untuk akses token.
type AccessClaims struct {
	Peran    []string `json:"peran"`
	HakAkses []string `json:"hak_akses"`
	jwt.RegisteredClaims
}

// TokenService bertugas membuat dan memverifikasi JWT serta refresh token.
type TokenService struct {
	cfg *config.Config
}

func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{cfg: cfg}
}

// GenerateAccessToken membuat JWT access token dengan daftar peran dan hak akses.
func (s *TokenService) GenerateAccessToken(userID uint, peran, hakAkses []string) (string, error) {
	now := time.Now()
	expires := now.Add(time.Duration(s.cfg.JWTAccessMinutes) * time.Minute)
	claims := AccessClaims{
		Peran:    peran,
		HakAkses: hakAkses,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTAccessSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// GenerateRefreshToken membuat token acak dan hash-nya untuk disimpan ke DB.
func (s *TokenService) GenerateRefreshToken() (plainToken string, hashed string, expires time.Time, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", time.Time{}, err
	}
	plainToken = base64.RawURLEncoding.EncodeToString(buf)
	hashed = s.HashRefreshToken(plainToken)
	expires = time.Now().Add(time.Duration(s.cfg.JWTRefreshDays) * 24 * time.Hour)
	return
}

// HashRefreshToken menghitung hash sha256 dari refresh token.
func (s *TokenService) HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// VerifyAccessToken memverifikasi JWT access token dan mengembalikan klaimnya.
func (s *TokenService) VerifyAccessToken(token string) (*AccessClaims, error) {
	parsed := &AccessClaims{}
	_, err := jwt.ParseWithClaims(token, parsed, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("metode signing tidak didukung")
		}
		return []byte(s.cfg.JWTAccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return parsed, nil
}
