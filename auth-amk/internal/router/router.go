package router

import (
	"amk-backend/auth-amk/internal/config"
	"amk-backend/auth-amk/internal/handler"
	"amk-backend/auth-amk/internal/middleware"
	"amk-backend/auth-amk/internal/repository"
	"amk-backend/auth-amk/internal/service"

	"github.com/gin-gonic/gin"
)

// NewRouter menyiapkan seluruh dependency dan route yang dibutuhkan auth-amk.
func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())

	// Repository & service
	penggunaRepo := repository.NewPenggunaRepository()
	peranRepo := repository.NewPeranRepository()
	tokenRepo := repository.NewTokenPenyegarRepository()
	resetRepo := repository.NewResetKataSandiRepository()

	tokenSvc := service.NewTokenService(cfg)
	authSvc := service.NewAuthService(penggunaRepo, peranRepo, tokenRepo, tokenSvc)
	hakSvc := service.NewHakAksesService(penggunaRepo)
	resetSvc := service.NewResetKataSandiService(penggunaRepo, resetRepo, tokenRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	lupaHandler := handler.NewLupaSandiHandler(resetSvc)
	meHandler := handler.NewMeHandler(penggunaRepo, hakSvc)

	// Health check
	r.GET("/health", handler.HealthHandler)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.POST("/lupa-sandi", lupaHandler.MintaOTP)
		authGroup.POST("/atur-sandi-baru", lupaHandler.ResetKataSandi)

		authGroup.GET("/me", middleware.JWTMiddleware(tokenSvc), meHandler.Me)
	}

	return r
}
