package model

// PenggunaPeran merepresentasikan tabel relasi pengguna_peran.
type PenggunaPeran struct {
	ID         uint `gorm:"primaryKey"`
	PenggunaID uint `gorm:"column:pengguna_id"`
	PeranID    uint `gorm:"column:peran_id"`
}

func (PenggunaPeran) TableName() string {
	return "pengguna_peran"
}
