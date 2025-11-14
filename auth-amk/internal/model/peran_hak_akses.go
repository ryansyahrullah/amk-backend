package model

// PeranHakAkses merepresentasikan tabel pivot peran_hak_akses.
type PeranHakAkses struct {
	ID         uint `gorm:"primaryKey"`
	PeranID    uint `gorm:"column:peran_id"`
	HakAksesID uint `gorm:"column:hak_akses_id"`
}

func (PeranHakAkses) TableName() string {
	return "peran_hak_akses"
}
