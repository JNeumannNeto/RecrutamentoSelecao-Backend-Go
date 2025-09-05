package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"recruitment-system/services/candidate-service/internal/application"
	"recruitment-system/services/candidate-service/internal/infrastructure"
	"recruitment-system/services/candidate-service/internal/interfaces"
	"recruitment-system/shared/database"
	"recruitment-system/shared/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	candidateRepo := infrastructure.NewCandidateRepository(db)
	skillRepo := infrastructure.NewSkillRepository(db)
	workExperienceRepo := infrastructure.NewWorkExperienceRepository(db)
	educationRepo := infrastructure.NewEducationRepository(db)
	resumeRepo := infrastructure.NewResumeRepository(db)
	jobApplicationRepo := infrastructure.NewJobApplicationRepository(db)

	authClient := infrastructure.NewAuthClient()
	jobClient := infrastructure.NewJobClient()
	fileStorageClient := infrastructure.NewFileStorageClient()
	aiClient := infrastructure.NewAIClient()

	candidateService := application.NewCandidateService(
		candidateRepo,
		skillRepo,
		workExperienceRepo,
		educationRepo,
		resumeRepo,
		jobApplicationRepo,
		authClient,
		jobClient,
		fileStorageClient,
		aiClient,
	)

	candidateController := interfaces.NewCandidateController(candidateService)

	r := gin.Default()

	r.Use(middleware.CORS())

	interfaces.SetupRoutes(r, candidateController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Candidate Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&infrastructure.CandidateModel{},
		&infrastructure.SkillModel{},
		&infrastructure.WorkExperienceModel{},
		&infrastructure.EducationModel{},
		&infrastructure.ResumeModel{},
		&infrastructure.JobApplicationModel{},
	)
}
