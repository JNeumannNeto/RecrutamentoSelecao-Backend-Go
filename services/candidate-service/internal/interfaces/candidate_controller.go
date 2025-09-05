package interfaces

import (
	"net/http"

	"recruitment-system/services/candidate-service/internal/application"
	"recruitment-system/services/candidate-service/internal/domain"
	"recruitment-system/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CandidateController struct {
	candidateService *application.CandidateService
}

func NewCandidateController(candidateService *application.CandidateService) *CandidateController {
	return &CandidateController{
		candidateService: candidateService,
	}
}

func (c *CandidateController) CreateCandidate(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	var req domain.CreateCandidateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	candidate, err := c.candidateService.CreateCandidate(ctx.Request.Context(), req, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create candidate profile", err)
		return
	}

	response := c.mapCandidateToResponse(candidate)
	utils.SuccessResponse(ctx, http.StatusCreated, "Candidate profile created successfully", response)
}

func (c *CandidateController) GetCandidate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	candidate, err := c.candidateService.GetCandidateByID(ctx.Request.Context(), id)
	if err != nil {
		utils.NotFoundResponse(ctx, "Candidate")
		return
	}

	response := c.mapCandidateToResponse(candidate)
	utils.SuccessResponse(ctx, http.StatusOK, "Candidate retrieved successfully", response)
}

func (c *CandidateController) GetMyProfile(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	candidate, err := c.candidateService.GetCandidateByUserID(ctx.Request.Context(), userInfo.ID)
	if err != nil {
		utils.NotFoundResponse(ctx, "Candidate profile")
		return
	}

	response := c.mapCandidateToResponse(candidate)
	utils.SuccessResponse(ctx, http.StatusOK, "Profile retrieved successfully", response)
}

func (c *CandidateController) UpdateCandidate(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	var req domain.UpdateCandidateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	candidate, err := c.candidateService.UpdateCandidate(ctx.Request.Context(), id, req, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update candidate profile", err)
		return
	}

	response := c.mapCandidateToResponse(candidate)
	utils.SuccessResponse(ctx, http.StatusOK, "Profile updated successfully", response)
}

func (c *CandidateController) AddSkill(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	candidateID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	var req domain.AddSkillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := c.candidateService.AddSkill(ctx.Request.Context(), candidateID, req, userInfo.ID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add skill", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Skill added successfully", nil)
}

func (c *CandidateController) RemoveSkill(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	candidateIDStr := ctx.Param("id")
	candidateID, err := uuid.Parse(candidateIDStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	skillIDStr := ctx.Param("skillId")
	skillID, err := uuid.Parse(skillIDStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid skill ID", err)
		return
	}

	if err := c.candidateService.RemoveSkill(ctx.Request.Context(), candidateID, skillID, userInfo.ID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove skill", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Skill removed successfully", nil)
}

func (c *CandidateController) AddWorkExperience(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	candidateID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	var req domain.AddWorkExperienceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	workExp, err := c.candidateService.AddWorkExperience(ctx.Request.Context(), candidateID, req, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add work experience", err)
		return
	}

	response := domain.WorkExperienceResponse{
		ID:          workExp.ID,
		CompanyName: workExp.CompanyName,
		Position:    workExp.Position,
		Description: workExp.Description,
		StartDate:   workExp.StartDate,
		EndDate:     workExp.EndDate,
		IsCurrent:   workExp.IsCurrent,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Work experience added successfully", response)
}

func (c *CandidateController) AddEducation(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	candidateID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	var req domain.AddEducationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	education, err := c.candidateService.AddEducation(ctx.Request.Context(), candidateID, req, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add education", err)
		return
	}

	response := domain.EducationResponse{
		ID:           education.ID,
		Institution:  education.Institution,
		Degree:       education.Degree,
		FieldOfStudy: education.FieldOfStudy,
		StartDate:    education.StartDate,
		EndDate:      education.EndDate,
		IsCurrent:    education.IsCurrent,
		GPA:          education.GPA,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Education added successfully", response)
}

func (c *CandidateController) UploadResume(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	candidateID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "No file uploaded", err)
		return
	}

	resume, err := c.candidateService.UploadResume(ctx.Request.Context(), candidateID, file, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to upload resume", err)
		return
	}

	response := domain.ResumeResponse{
		ID:          resume.ID,
		Filename:    resume.Filename,
		FileSize:    resume.FileSize,
		MimeType:    resume.MimeType,
		AIProcessed: resume.AIProcessed,
		UploadedAt:  resume.CreatedAt,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Resume uploaded successfully", response)
}

func (c *CandidateController) ApplyToJob(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	candidateID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	var req domain.CreateJobApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	application, err := c.candidateService.ApplyToJob(ctx.Request.Context(), candidateID, req, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to apply to job", err)
		return
	}

	response := domain.JobApplicationResponse{
		ID:          application.ID,
		Status:      application.Status,
		CoverLetter: application.CoverLetter,
		AppliedAt:   application.AppliedAt,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Application submitted successfully", response)
}

func (c *CandidateController) GetApplications(ctx *gin.Context) {
	token := c.extractToken(ctx)
	if token == "" {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userInfo, err := c.candidateService.ValidateUserPermissions(ctx.Request.Context(), token, "candidate")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	idStr := ctx.Param("id")
	candidateID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid candidate ID", err)
		return
	}

	applications, err := c.candidateService.GetApplications(ctx.Request.Context(), candidateID, userInfo.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to get applications", err)
		return
	}

	responses := make([]domain.JobApplicationResponse, len(applications))
	for i, app := range applications {
		responses[i] = domain.JobApplicationResponse{
			ID:          app.ID,
			Status:      app.Status,
			CoverLetter: app.CoverLetter,
			AppliedAt:   app.AppliedAt,
		}
		if app.Job != nil {
			responses[i].Job = domain.JobResponse{
				ID:          app.Job.ID,
				Title:       app.Job.Title,
				Description: app.Job.Description,
				Location:    app.Job.Location,
				Status:      app.Job.Status,
			}
		}
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Applications retrieved successfully", responses)
}

func (c *CandidateController) extractToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return ""
}

func (c *CandidateController) mapCandidateToResponse(candidate *domain.Candidate) domain.CandidateResponse {
	response := domain.CandidateResponse{
		ID:          candidate.ID,
		Phone:       candidate.Phone,
		Address:     candidate.Address,
		DateOfBirth: candidate.DateOfBirth,
		LinkedinURL: candidate.LinkedinURL,
		GithubURL:   candidate.GithubURL,
		CreatedAt:   candidate.CreatedAt,
	}

	if candidate.User != nil {
		response.User = domain.UserResponse{
			ID:    candidate.User.ID,
			Email: candidate.User.Email,
			Name:  candidate.User.Name,
			Role:  candidate.User.Role,
		}
	}

	if len(candidate.Skills) > 0 {
		response.Skills = make([]domain.CandidateSkillResponse, len(candidate.Skills))
		for i, skill := range candidate.Skills {
			response.Skills[i] = domain.CandidateSkillResponse{
				ID:                skill.ID,
				ProficiencyLevel:  skill.ProficiencyLevel,
				YearsOfExperience: skill.YearsOfExperience,
			}
			if skill.Skill != nil {
				response.Skills[i].Skill = domain.SkillResponse{
					ID:       skill.Skill.ID,
					Name:     skill.Skill.Name,
					Category: skill.Skill.Category,
				}
			}
		}
	}

	if len(candidate.WorkExperiences) > 0 {
		response.WorkExperiences = make([]domain.WorkExperienceResponse, len(candidate.WorkExperiences))
		for i, exp := range candidate.WorkExperiences {
			response.WorkExperiences[i] = domain.WorkExperienceResponse{
				ID:          exp.ID,
				CompanyName: exp.CompanyName,
				Position:    exp.Position,
				Description: exp.Description,
				StartDate:   exp.StartDate,
				EndDate:     exp.EndDate,
				IsCurrent:   exp.IsCurrent,
			}
		}
	}

	if len(candidate.Education) > 0 {
		response.Education = make([]domain.EducationResponse, len(candidate.Education))
		for i, edu := range candidate.Education {
			response.Education[i] = domain.EducationResponse{
				ID:           edu.ID,
				Institution:  edu.Institution,
				Degree:       edu.Degree,
				FieldOfStudy: edu.FieldOfStudy,
				StartDate:    edu.StartDate,
				EndDate:      edu.EndDate,
				IsCurrent:    edu.IsCurrent,
				GPA:          edu.GPA,
			}
		}
	}

	if len(candidate.Resumes) > 0 {
		response.Resumes = make([]domain.ResumeResponse, len(candidate.Resumes))
		for i, resume := range candidate.Resumes {
			response.Resumes[i] = domain.ResumeResponse{
				ID:          resume.ID,
				Filename:    resume.Filename,
				FileSize:    resume.FileSize,
				MimeType:    resume.MimeType,
				AIProcessed: resume.AIProcessed,
				UploadedAt:  resume.CreatedAt,
			}
		}
	}

	return response
}
