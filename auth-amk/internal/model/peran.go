package model

import "time"

// Peran merepresentasikan role logis di sistem.
type Peran struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Nama       string    `gorm:"size:100;uniqueIndex" json:"nama"`
	Keterangan string    `gorm:"size:255" json:"keterangan"`
	DibuatPada time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiubahPada time.Time `gorm:"column:diubah_pada;autoUpdateTime" json:"diubah_pada"`

	HakAkses []*HakAkses `gorm:"many2many:peran_hak_akses" json:"hak_akses,omitempty"`
}

func (Peran) TableName() string {
	return "peran"
}
