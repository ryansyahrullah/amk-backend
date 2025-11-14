// auth-amk/cmd/auth/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"amk-backend/auth-amk/internal/config"
	"amk-backend/auth-amk/internal/database"
	"amk-backend/auth-amk/internal/router"
)

func main() {
	// 1. Load konfigurasi dari .env
	cfg := config.LoadConfig()

	// 2. Koneksi ke database auth_db
	if err := database.ConnectDatabase(cfg); err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}

	// 3. Jalankan migrasi tabel auth (pengguna, peran, dll)
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("gagal migrate database: %v", err)
	}

	// 4. Inisialisasi router Gin (semua route & middleware didaftarkan di sini)
	engine := router.NewRouter(cfg)

	// 5. Setup HTTP server (supaya bisa graceful shutdown)
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: engine,
	}

	// 6. Jalankan server di goroutine terpisah
	go func() {
		log.Printf("[AUTH] berjalan di port %s, env=%s\n", cfg.AppPort, cfg.AppEnv)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 7. Tunggu sinyal interrupt (Ctrl+C atau SIGTERM di server)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[AUTH] menerima sinyal shutdown, mematikan server...")

	// 8. Graceful shutdown dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[AUTH] shutdown paksa: %v", err)
	}

	log.Println("[AUTH] server berhenti dengan aman")
}
