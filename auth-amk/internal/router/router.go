package router

import (
	"github.com/gin-gonic/gin"

	"amk-backend/auth-amk/internal/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Endpoint simple untuk cek service hidup
	r.GET("/health", handler.HealthHandler)

	// Nanti di sini kita daftarkan:
	// - POST /auth/login
	// - POST /auth/refresh
	// - POST /auth/logout
	// - POST /auth/lupa-sandi
	// - POST /auth/atur-sandi-baru
	// - GET  /auth/me
	// dengan kontrol role + hak_akses sesuai desain.

	return r
}
