package interfaces

import (
	"recruitment-system/shared/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authController *AuthController, jwtSecret string) {
	api := router.Group("/api/v1")
	
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.RefreshToken)
		auth.POST("/validate", authController.ValidateToken)
	}

	protected := api.Group("/auth")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.POST("/logout", authController.Logout)
		protected.GET("/profile", authController.GetProfile)
		protected.PUT("/change-password", authController.ChangePassword)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "auth-service",
		})
	})
}
