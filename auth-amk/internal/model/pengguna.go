package model

import "time"

// Pengguna merepresentasikan akun yang dapat mengakses sistem otentikasi.
type Pengguna struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Email         string    `gorm:"size:191;uniqueIndex" json:"email"`
	NRP           string    `gorm:"size:50;uniqueIndex" json:"nrp"`
	KataSandiHash string    `gorm:"column:kata_sandi_hash" json:"-"`
	Status        string    `gorm:"size:20" json:"status"`
	DibuatPada    time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiubahPada    time.Time `gorm:"column:diubah_pada;autoUpdateTime" json:"diubah_pada"`

	Peran []*Peran `gorm:"many2many:pengguna_peran" json:"peran,omitempty"`
}

// TableName menyesuaikan nama tabel sebenarnya di database.
func (Pengguna) TableName() string {
	return "pengguna"
}

// Aktif mengembalikan true jika status pengguna masih aktif.
func (p Pengguna) Aktif() bool {
	return p.Status == "aktif"
}
