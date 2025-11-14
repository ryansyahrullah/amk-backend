package model

import "time"

// ResetKataSandi menyimpan kode OTP untuk proses lupa kata sandi.
type ResetKataSandi struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PenggunaID      uint      `gorm:"column:pengguna_id" json:"pengguna_id"`
	KodeOTPHash     string    `gorm:"column:kode_otp_hash;size:255" json:"-"`
	BerlakuSampai   time.Time `gorm:"column:berlaku_sampai" json:"berlaku_sampai"`
	JumlahPercobaan int       `gorm:"column:jumlah_percobaan" json:"jumlah_percobaan"`
	Status          string    `gorm:"size:20" json:"status"`
	DibuatPada      time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
}

func (ResetKataSandi) TableName() string {
	return "reset_kata_sandi"
}
