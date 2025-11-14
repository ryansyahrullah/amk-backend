package model

import "time"

// HakAkses adalah entitas permission granular.
type HakAkses struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Kode       string    `gorm:"size:191;uniqueIndex" json:"kode"`
	Nama       string    `gorm:"size:191" json:"nama"`
	Keterangan string    `gorm:"size:255" json:"keterangan"`
	DibuatPada time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiubahPada time.Time `gorm:"column:diubah_pada;autoUpdateTime" json:"diubah_pada"`
}

func (HakAkses) TableName() string {
	return "hak_akses"
}
