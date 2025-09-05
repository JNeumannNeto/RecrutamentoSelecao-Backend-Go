package interfaces

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, candidateController *CandidateController) {
	api := router.Group("/api/v1")

	candidates := api.Group("/candidates")
	{
		candidates.POST("", candidateController.CreateCandidate)
		candidates.GET("/profile", candidateController.GetMyProfile)
		candidates.GET("/:id", candidateController.GetCandidate)
		candidates.PUT("/:id", candidateController.UpdateCandidate)
		
		candidates.POST("/:id/skills", candidateController.AddSkill)
		candidates.DELETE("/:id/skills/:skillId", candidateController.RemoveSkill)
		
		candidates.POST("/:id/work-experiences", candidateController.AddWorkExperience)
		candidates.POST("/:id/education", candidateController.AddEducation)
		
		candidates.POST("/:id/resume", candidateController.UploadResume)
		
		candidates.POST("/:id/applications", candidateController.ApplyToJob)
		candidates.GET("/:id/applications", candidateController.GetApplications)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "candidate-service",
		})
	})
}
