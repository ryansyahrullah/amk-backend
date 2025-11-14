package model

import "time"

// TokenPenyegar merepresentasikan refresh token yang disimpan sebagai hash.
type TokenPenyegar struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	PenggunaID    uint       `gorm:"column:pengguna_id" json:"pengguna_id"`
	TokenHash     string     `gorm:"column:token_hash;size:255" json:"-"`
	UserAgent     string     `gorm:"column:user_agent;size:255" json:"user_agent"`
	AlamatIP      string     `gorm:"column:alamat_ip;size:100" json:"alamat_ip"`
	BerlakuSampai time.Time  `gorm:"column:berlaku_sampai" json:"berlaku_sampai"`
	DicabutPada   *time.Time `gorm:"column:dicabut_pada" json:"dicabut_pada"`
	DibuatPada    time.Time  `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
}

func (TokenPenyegar) TableName() string {
	return "token_penyegar"
}

// Aktif mengembalikan true jika token belum dicabut dan belum kedaluwarsa.
func (t TokenPenyegar) Aktif(now time.Time) bool {
	if t.DicabutPada != nil {
		return false
	}
	return now.Before(t.BerlakuSampai)
}
