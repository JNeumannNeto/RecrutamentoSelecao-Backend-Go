package interfaces

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, jobController *JobController, skillController *SkillController) {
	api := router.Group("/api/v1")

	jobs := api.Group("/jobs")
	{
		jobs.GET("", jobController.ListJobs)
		jobs.GET("/:id", jobController.GetJob)
		jobs.POST("", jobController.CreateJob)
		jobs.PUT("/:id", jobController.UpdateJob)
		jobs.DELETE("/:id", jobController.DeleteJob)
		jobs.PATCH("/:id/status", jobController.UpdateJobStatus)
		jobs.GET("/my", jobController.GetMyJobs)
	}

	skills := api.Group("/skills")
	{
		skills.GET("", skillController.ListSkills)
		skills.POST("", skillController.CreateSkill)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "job-service",
		})
	})
}
