package database

import (
	"fmt"
	"log"
	"time"

	"amk-backend/auth-amk/internal/config"
	"amk-backend/auth-amk/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var logLevel logger.LogLevel
	switch cfg.DBLogMode {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	default:
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("gagal konek ke DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("gagal ambil sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	DB = db
	log.Println("[AUTH] koneksi ke database sukses")
	return nil
}

// AutoMigrate akan memastikan seluruh tabel penting tersedia di database.
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database belum terhubung")
	}

	return DB.AutoMigrate(
		&model.Pengguna{},
		&model.Peran{},
		&model.PenggunaPeran{},
		&model.HakAkses{},
		&model.PeranHakAkses{},
		&model.TokenPenyegar{},
		&model.ResetKataSandi{},
	)
}
