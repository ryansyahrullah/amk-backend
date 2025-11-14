package middleware

import (
	"net/http"
	"strconv"

	"amk-backend/auth-amk/internal/service"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware memvalidasi access token dan meletakkan klaim ke context.
func JWTMiddleware(tokenSvc *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token tidak ditemukan"})
			return
		}
		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "format token tidak valid"})
			return
		}

		token := authHeader[len(prefix):]
		claims, err := tokenSvc.VerifyAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token tidak valid"})
			return
		}

		userID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token tidak valid"})
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", uint(userID))
		c.Next()
	}
}
