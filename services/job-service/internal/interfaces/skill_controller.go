package interfaces

import (
	"net/http"

	"recruitment-system/services/job-service/internal/application"
	"recruitment-system/services/job-service/internal/domain"
	"recruitment-system/shared/utils"

	"github.com/gin-gonic/gin"
)

type SkillController struct {
	jobService *application.JobService
}

func NewSkillController(jobService *application.JobService) *SkillController {
	return &SkillController{
		jobService: jobService,
	}
}

func (c *SkillController) ListSkills(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	category := ctx.Query("category")
	search := ctx.Query("search")

	skills, total, err := c.jobService.ListSkills(ctx.Request.Context(), category, search, pagination.Page, pagination.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err)
		return
	}

	responses := make([]domain.SkillResponse, len(skills))
	for i, skill := range skills {
		responses[i] = domain.SkillResponse{
			ID:       skill.ID,
			Name:     skill.Name,
			Category: skill.Category,
		}
	}

	paginationInfo := utils.CreatePagination(pagination.Page, pagination.Limit, total)
	utils.PaginatedSuccessResponse(ctx, http.StatusOK, "Skills retrieved successfully", responses, paginationInfo)
}

func (c *SkillController) CreateSkill(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	_, err := c.jobService.ValidateUserPermissions(ctx.Request.Context(), token, "admin")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required"`
		Category string `json:"category"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	skill, err := c.jobService.CreateSkill(ctx.Request.Context(), req.Name, req.Category)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create skill", err)
		return
	}

	response := domain.SkillResponse{
		ID:       skill.ID,
		Name:     skill.Name,
		Category: skill.Category,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Skill created successfully", response)
}

func (c *SkillController) extractToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return ""
}
