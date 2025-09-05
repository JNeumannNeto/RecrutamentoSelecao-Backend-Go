package main

import (
	"log"
	"os"

	"recruitment-system/services/job-service/internal/application"
	"recruitment-system/services/job-service/internal/infrastructure"
	"recruitment-system/services/job-service/internal/interfaces"
	"recruitment-system/shared/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbConfig := database.GetConfigFromEnv()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.TestConnection(db); err != nil {
		log.Fatal("Database connection test failed:", err)
	}

	jobRepo := infrastructure.NewJobRepository(db)
	skillRepo := infrastructure.NewSkillRepository(db)
	jobSkillRepo := infrastructure.NewJobSkillRepository(db)

	authServiceURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8083")
	authClient := infrastructure.NewAuthServiceClient(authServiceURL)

	jobService := application.NewJobService(jobRepo, skillRepo, jobSkillRepo, authClient)

	jobController := interfaces.NewJobController(jobService)
	skillController := interfaces.NewSkillController(jobService)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	interfaces.SetupRoutes(router, jobController, skillController)

	port := getEnv("PORT", "8081")
	log.Printf("Job Service starting on port %s", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
