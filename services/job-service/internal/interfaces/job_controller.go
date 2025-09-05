package interfaces

import (
	"net/http"
	"strconv"

	"recruitment-system/services/job-service/internal/application"
	"recruitment-system/services/job-service/internal/domain"
	"recruitment-system/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobController struct {
	jobService *application.JobService
}

func NewJobController(jobService *application.JobService) *JobController {
	return &JobController{
		jobService: jobService,
	}
}

func (c *JobController) CreateJob(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.jobService.ValidateUserPermissions(ctx.Request.Context(), token, "admin")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	var req domain.CreateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	job, err := c.jobService.CreateJob(ctx.Request.Context(), req, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create job", err)
		return
	}

	response := c.mapJobToResponse(job)
	utils.SuccessResponse(ctx, http.StatusCreated, "Job created successfully", response)
}

func (c *JobController) GetJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid job ID", err)
		return
	}

	job, err := c.jobService.GetJobByID(ctx.Request.Context(), id)
	if err != nil {
		utils.NotFoundResponse(ctx, "Job")
		return
	}

	response := c.mapJobToResponse(job)
	utils.SuccessResponse(ctx, http.StatusOK, "Job retrieved successfully", response)
}

func (c *JobController) UpdateJob(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.jobService.ValidateUserPermissions(ctx.Request.Context(), token, "admin")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid job ID", err)
		return
	}

	var req domain.UpdateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	job, err := c.jobService.UpdateJob(ctx.Request.Context(), id, req, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update job", err)
		return
	}

	response := c.mapJobToResponse(job)
	utils.SuccessResponse(ctx, http.StatusOK, "Job updated successfully", response)
}

func (c *JobController) UpdateJobStatus(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.jobService.ValidateUserPermissions(ctx.Request.Context(), token, "admin")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid job ID", err)
		return
	}

	var req domain.UpdateJobStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := c.jobService.UpdateJobStatus(ctx.Request.Context(), id, req.Status, userInfo.ID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update job status", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Job status updated successfully", nil)
}

func (c *JobController) DeleteJob(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.jobService.ValidateUserPermissions(ctx.Request.Context(), token, "admin")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid job ID", err)
		return
	}

	if err := c.jobService.DeleteJob(ctx.Request.Context(), id, userInfo.ID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete job", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Job deleted successfully", nil)
}

func (c *JobController) ListJobs(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	filter := domain.JobListFilter{
		Status:   ctx.Query("status"),
		Location: ctx.Query("location"),
		Title:    ctx.Query("title"),
	}

	if minSalaryStr := ctx.Query("min_salary"); minSalaryStr != "" {
		if minSalary, err := strconv.ParseFloat(minSalaryStr, 64); err == nil {
			filter.MinSalary = &minSalary
		}
	}

	if maxSalaryStr := ctx.Query("max_salary"); maxSalaryStr != "" {
		if maxSalary, err := strconv.ParseFloat(maxSalaryStr, 64); err == nil {
			filter.MaxSalary = &maxSalary
		}
	}

	jobs, total, err := c.jobService.ListJobs(ctx.Request.Context(), filter, pagination.Page, pagination.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err)
		return
	}

	responses := make([]domain.JobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = c.mapJobToResponse(job)
	}

	paginationInfo := utils.CreatePagination(pagination.Page, pagination.Limit, total)
	utils.PaginatedSuccessResponse(ctx, http.StatusOK, "Jobs retrieved successfully", responses, paginationInfo)
}

func (c *JobController) GetMyJobs(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.jobService.ValidateUserPermissions(ctx.Request.Context(), token, "admin")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	pagination := utils.GetPaginationParams(ctx)

	jobs, total, err := c.jobService.GetJobsByCreatedBy(ctx.Request.Context(), userInfo.ID, pagination.Page, pagination.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err)
		return
	}

	responses := make([]domain.JobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = c.mapJobToResponse(job)
	}

	paginationInfo := utils.CreatePagination(pagination.Page, pagination.Limit, total)
	utils.PaginatedSuccessResponse(ctx, http.StatusOK, "Jobs retrieved successfully", responses, paginationInfo)
}

func (c *JobController) extractToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return ""
}

func (c *JobController) mapJobToResponse(job *domain.Job) domain.JobResponse {
	response := domain.JobResponse{
		ID:           job.ID,
		Title:        job.Title,
		Description:  job.Description,
		Requirements: job.Requirements,
		Location:     job.Location,
		SalaryMin:    job.SalaryMin,
		SalaryMax:    job.SalaryMax,
		Status:       job.Status,
		CreatedBy:    job.CreatedBy,
		CreatedAt:    job.CreatedAt,
		UpdatedAt:    job.UpdatedAt,
	}

	if len(job.Skills) > 0 {
		response.Skills = make([]domain.JobSkillResponse, len(job.Skills))
		for i, jobSkill := range job.Skills {
			skillResponse := domain.SkillResponse{
				ID:       jobSkill.Skill.ID,
				Name:     jobSkill.Skill.Name,
				Category: jobSkill.Skill.Category,
			}
			if jobSkill.Skill != nil {
				skillResponse = domain.SkillResponse{
					ID:       jobSkill.Skill.ID,
					Name:     jobSkill.Skill.Name,
					Category: jobSkill.Skill.Category,
				}
			}

			response.Skills[i] = domain.JobSkillResponse{
				Skill:         skillResponse,
				RequiredLevel: jobSkill.RequiredLevel,
				IsRequired:    jobSkill.IsRequired,
			}
		}
	}

	return response
}
