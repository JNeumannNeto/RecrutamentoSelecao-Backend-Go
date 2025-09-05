package application

import (
	"context"
	"errors"
	"time"

	"recruitment-system/services/job-service/internal/domain"
	"recruitment-system/shared/utils"

	"github.com/google/uuid"
)

type JobService struct {
	jobRepo      domain.JobRepository
	skillRepo    domain.SkillRepository
	jobSkillRepo domain.JobSkillRepository
	authClient   domain.AuthServiceClient
}

func NewJobService(
	jobRepo domain.JobRepository,
	skillRepo domain.SkillRepository,
	jobSkillRepo domain.JobSkillRepository,
	authClient domain.AuthServiceClient,
) *JobService {
	return &JobService{
		jobRepo:      jobRepo,
		skillRepo:    skillRepo,
		jobSkillRepo: jobSkillRepo,
		authClient:   authClient,
	}
}

func (s *JobService) CreateJob(ctx context.Context, req domain.CreateJobRequest, createdBy uuid.UUID) (*domain.Job, error) {
	if utils.IsEmptyOrWhitespace(req.Title) {
		return nil, errors.New("title is required")
	}

	if utils.IsEmptyOrWhitespace(req.Description) {
		return nil, errors.New("description is required")
	}

	if req.SalaryMin != nil && req.SalaryMax != nil && *req.SalaryMin > *req.SalaryMax {
		return nil, errors.New("minimum salary cannot be greater than maximum salary")
	}

	if len(req.Skills) > 0 {
		skillIDs := make([]uuid.UUID, len(req.Skills))
		for i, skill := range req.Skills {
			skillIDs[i] = skill.SkillID
		}

		skills, err := s.skillRepo.GetByIDs(ctx, skillIDs)
		if err != nil {
			return nil, err
		}

		if len(skills) != len(skillIDs) {
			return nil, errors.New("one or more skills not found")
		}
	}

	job := &domain.Job{
		ID:           uuid.New(),
		Title:        utils.SanitizeString(req.Title),
		Description:  utils.SanitizeString(req.Description),
		Requirements: utils.SanitizeString(req.Requirements),
		Location:     utils.SanitizeString(req.Location),
		SalaryMin:    req.SalaryMin,
		SalaryMax:    req.SalaryMax,
		Status:       string(domain.JobStatusOpen),
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, err
	}

	if len(req.Skills) > 0 {
		jobSkills := make([]domain.JobSkill, len(req.Skills))
		for i, skillReq := range req.Skills {
			jobSkills[i] = domain.JobSkill{
				ID:            uuid.New(),
				JobID:         job.ID,
				SkillID:       skillReq.SkillID,
				RequiredLevel: skillReq.RequiredLevel,
				IsRequired:    skillReq.IsRequired,
				CreatedAt:     time.Now(),
			}
		}

		if err := s.jobSkillRepo.CreateBatch(ctx, jobSkills); err != nil {
			return nil, err
		}
	}

	return s.GetJobByID(ctx, job.ID)
}

func (s *JobService) GetJobByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	jobSkills, err := s.jobSkillRepo.GetByJobID(ctx, job.ID)
	if err != nil {
		return nil, err
	}

	job.Skills = jobSkills
	return job, nil
}

func (s *JobService) UpdateJob(ctx context.Context, id uuid.UUID, req domain.UpdateJobRequest, userID uuid.UUID) (*domain.Job, error) {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if job.CreatedBy != userID {
		return nil, errors.New("you can only update jobs you created")
	}

	if req.Title != "" {
		job.Title = utils.SanitizeString(req.Title)
	}
	if req.Description != "" {
		job.Description = utils.SanitizeString(req.Description)
	}
	if req.Requirements != "" {
		job.Requirements = utils.SanitizeString(req.Requirements)
	}
	if req.Location != "" {
		job.Location = utils.SanitizeString(req.Location)
	}
	if req.SalaryMin != nil {
		job.SalaryMin = req.SalaryMin
	}
	if req.SalaryMax != nil {
		job.SalaryMax = req.SalaryMax
	}

	if job.SalaryMin != nil && job.SalaryMax != nil && *job.SalaryMin > *job.SalaryMax {
		return nil, errors.New("minimum salary cannot be greater than maximum salary")
	}

	job.UpdatedAt = time.Now()

	if err := s.jobRepo.Update(ctx, job); err != nil {
		return nil, err
	}

	return s.GetJobByID(ctx, job.ID)
}

func (s *JobService) UpdateJobStatus(ctx context.Context, id uuid.UUID, status string, userID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if job.CreatedBy != userID {
		return errors.New("you can only update jobs you created")
	}

	if !utils.IsValidJobStatus(status) {
		return errors.New("invalid job status")
	}

	return s.jobRepo.UpdateStatus(ctx, id, status)
}

func (s *JobService) DeleteJob(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if job.CreatedBy != userID {
		return errors.New("you can only delete jobs you created")
	}

	if err := s.jobSkillRepo.DeleteByJobID(ctx, id); err != nil {
		return err
	}

	return s.jobRepo.Delete(ctx, id)
}

func (s *JobService) ListJobs(ctx context.Context, filter domain.JobListFilter, page, limit int) ([]*domain.Job, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := utils.CalculateOffset(page, limit)
	jobs, total, err := s.jobRepo.List(ctx, filter, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for _, job := range jobs {
		jobSkills, err := s.jobSkillRepo.GetByJobID(ctx, job.ID)
		if err != nil {
			continue
		}
		job.Skills = jobSkills
	}

	return jobs, total, nil
}

func (s *JobService) GetJobsByCreatedBy(ctx context.Context, createdBy uuid.UUID, page, limit int) ([]*domain.Job, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := utils.CalculateOffset(page, limit)
	jobs, total, err := s.jobRepo.GetByCreatedBy(ctx, createdBy, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for _, job := range jobs {
		jobSkills, err := s.jobSkillRepo.GetByJobID(ctx, job.ID)
		if err != nil {
			continue
		}
		job.Skills = jobSkills
	}

	return jobs, total, nil
}

func (s *JobService) ListSkills(ctx context.Context, category, search string, page, limit int) ([]*domain.Skill, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := utils.CalculateOffset(page, limit)
	return s.skillRepo.List(ctx, category, search, offset, limit)
}

func (s *JobService) CreateSkill(ctx context.Context, name, category string) (*domain.Skill, error) {
	if utils.IsEmptyOrWhitespace(name) {
		return nil, errors.New("skill name is required")
	}

	exists, err := s.skillRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("skill with this name already exists")
	}

	skill := &domain.Skill{
		ID:        uuid.New(),
		Name:      utils.SanitizeString(name),
		Category:  utils.SanitizeString(category),
		CreatedAt: time.Now(),
	}

	if err := s.skillRepo.Create(ctx, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *JobService) ValidateUserPermissions(ctx context.Context, token string, requiredRole string) (*domain.UserInfo, error) {
	userInfo, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if requiredRole != "" && userInfo.Role != requiredRole {
		return nil, errors.New("insufficient permissions")
	}

	return userInfo, nil
}
