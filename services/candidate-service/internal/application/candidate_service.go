package application

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"recruitment-system/services/candidate-service/internal/domain"
	"recruitment-system/shared/utils"

	"github.com/google/uuid"
)

type CandidateService struct {
	candidateRepo     domain.CandidateRepository
	candidateSkillRepo domain.CandidateSkillRepository
	workExpRepo       domain.WorkExperienceRepository
	educationRepo     domain.EducationRepository
	resumeRepo        domain.ResumeRepository
	applicationRepo   domain.JobApplicationRepository
	skillRepo         domain.SkillRepository
	fileStorage       domain.FileStorageService
	aiService         domain.AIService
	authClient        domain.AuthServiceClient
	jobClient         domain.JobServiceClient
}

func NewCandidateService(
	candidateRepo domain.CandidateRepository,
	candidateSkillRepo domain.CandidateSkillRepository,
	workExpRepo domain.WorkExperienceRepository,
	educationRepo domain.EducationRepository,
	resumeRepo domain.ResumeRepository,
	applicationRepo domain.JobApplicationRepository,
	skillRepo domain.SkillRepository,
	fileStorage domain.FileStorageService,
	aiService domain.AIService,
	authClient domain.AuthServiceClient,
	jobClient domain.JobServiceClient,
) *CandidateService {
	return &CandidateService{
		candidateRepo:      candidateRepo,
		candidateSkillRepo: candidateSkillRepo,
		workExpRepo:        workExpRepo,
		educationRepo:      educationRepo,
		resumeRepo:         resumeRepo,
		applicationRepo:    applicationRepo,
		skillRepo:          skillRepo,
		fileStorage:        fileStorage,
		aiService:          aiService,
		authClient:         authClient,
		jobClient:          jobClient,
	}
}

func (s *CandidateService) CreateCandidate(ctx context.Context, req domain.CreateCandidateRequest, userID uuid.UUID) (*domain.Candidate, error) {
	exists, err := s.candidateRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("candidate profile already exists for this user")
	}

	candidate := &domain.Candidate{
		ID:          uuid.New(),
		UserID:      userID,
		Phone:       utils.SanitizeString(req.Phone),
		Address:     utils.SanitizeString(req.Address),
		DateOfBirth: req.DateOfBirth,
		LinkedinURL: utils.SanitizeString(req.LinkedinURL),
		GithubURL:   utils.SanitizeString(req.GithubURL),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.candidateRepo.Create(ctx, candidate); err != nil {
		return nil, err
	}

	return s.GetCandidateByID(ctx, candidate.ID)
}

func (s *CandidateService) GetCandidateByID(ctx context.Context, id uuid.UUID) (*domain.Candidate, error) {
	candidate, err := s.candidateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.loadCandidateRelations(ctx, candidate); err != nil {
		return nil, err
	}

	return candidate, nil
}

func (s *CandidateService) GetCandidateByUserID(ctx context.Context, userID uuid.UUID) (*domain.Candidate, error) {
	candidate, err := s.candidateRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.loadCandidateRelations(ctx, candidate); err != nil {
		return nil, err
	}

	return candidate, nil
}

func (s *CandidateService) UpdateCandidate(ctx context.Context, id uuid.UUID, req domain.UpdateCandidateRequest, userID uuid.UUID) (*domain.Candidate, error) {
	candidate, err := s.candidateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if candidate.UserID != userID {
		return nil, errors.New("you can only update your own profile")
	}

	candidate.Phone = utils.SanitizeString(req.Phone)
	candidate.Address = utils.SanitizeString(req.Address)
	candidate.DateOfBirth = req.DateOfBirth
	candidate.LinkedinURL = utils.SanitizeString(req.LinkedinURL)
	candidate.GithubURL = utils.SanitizeString(req.GithubURL)
	candidate.UpdatedAt = time.Now()

	if err := s.candidateRepo.Update(ctx, candidate); err != nil {
		return nil, err
	}

	return s.GetCandidateByID(ctx, candidate.ID)
}

func (s *CandidateService) AddSkill(ctx context.Context, candidateID uuid.UUID, req domain.AddSkillRequest, userID uuid.UUID) error {
	candidate, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil {
		return err
	}

	if candidate.UserID != userID {
		return errors.New("you can only add skills to your own profile")
	}

	skill, err := s.skillRepo.GetByID(ctx, req.SkillID)
	if err != nil {
		return errors.New("skill not found")
	}

	if !utils.IsValidProficiencyLevel(req.ProficiencyLevel) {
		return errors.New("invalid proficiency level")
	}

	candidateSkill := &domain.CandidateSkill{
		ID:                uuid.New(),
		CandidateID:       candidateID,
		SkillID:           req.SkillID,
		ProficiencyLevel:  req.ProficiencyLevel,
		YearsOfExperience: req.YearsOfExperience,
		CreatedAt:         time.Now(),
	}

	return s.candidateSkillRepo.Create(ctx, candidateSkill)
}

func (s *CandidateService) RemoveSkill(ctx context.Context, candidateID, skillID uuid.UUID, userID uuid.UUID) error {
	candidate, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil {
		return err
	}

	if candidate.UserID != userID {
		return errors.New("you can only remove skills from your own profile")
	}

	return s.candidateSkillRepo.DeleteByCandidateIDAndSkillID(ctx, candidateID, skillID)
}

func (s *CandidateService) AddWorkExperience(ctx context.Context, candidateID uuid.UUID, req domain.AddWorkExperienceRequest, userID uuid.UUID) (*domain.WorkExperience, error) {
	candidate, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil {
		return nil, err
	}

	if candidate.UserID != userID {
		return nil, errors.New("you can only add work experience to your own profile")
	}

	workExp := &domain.WorkExperience{
		ID:          uuid.New(),
		CandidateID: candidateID,
		CompanyName: utils.SanitizeString(req.CompanyName),
		Position:    utils.SanitizeString(req.Position),
		Description: utils.SanitizeString(req.Description),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsCurrent:   req.IsCurrent,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.workExpRepo.Create(ctx, workExp); err != nil {
		return nil, err
	}

	return workExp, nil
}

func (s *CandidateService) AddEducation(ctx context.Context, candidateID uuid.UUID, req domain.AddEducationRequest, userID uuid.UUID) (*domain.Education, error) {
	candidate, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil {
		return nil, err
	}

	if candidate.UserID != userID {
		return nil, errors.New("you can only add education to your own profile")
	}

	education := &domain.Education{
		ID:           uuid.New(),
		CandidateID:  candidateID,
		Institution:  utils.SanitizeString(req.Institution),
		Degree:       utils.SanitizeString(req.Degree),
		FieldOfStudy: utils.SanitizeString(req.FieldOfStudy),
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		IsCurrent:    req.IsCurrent,
		GPA:          req.GPA,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.educationRepo.Create(ctx, education); err != nil {
		return nil, err
	}

	return education, nil
}

func (s *CandidateService) UploadResume(ctx context.Context, candidateID uuid.UUID, file *multipart.FileHeader, userID uuid.UUID) (*domain.Resume, error) {
	candidate, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil {
		return nil, err
	}

	if candidate.UserID != userID {
		return nil, errors.New("you can only upload resume to your own profile")
	}

	filePath, err := s.fileStorage.SaveFile(ctx, file, candidateID)
	if err != nil {
		return nil, err
	}

	resume := &domain.Resume{
		ID:          uuid.New(),
		CandidateID: candidateID,
		Filename:    file.Filename,
		FilePath:    filePath,
		FileSize:    file.Size,
		MimeType:    file.Header.Get("Content-Type"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.resumeRepo.Create(ctx, resume); err != nil {
		s.fileStorage.DeleteFile(ctx, filePath)
		return nil, err
	}

	go s.processResumeWithAI(context.Background(), resume)

	return resume, nil
}

func (s *CandidateService) ApplyToJob(ctx context.Context, candidateID uuid.UUID, req domain.CreateJobApplicationRequest, userID uuid.UUID) (*domain.JobApplication, error) {
	candidate, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil {
		return nil, err
	}

	if candidate.UserID != userID {
		return nil, errors.New("you can only apply to jobs with your own profile")
	}

	exists, err := s.applicationRepo.ExistsByCandidateAndJob(ctx, candidateID, req.JobID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("you have already applied to this job")
	}

	isOpen, err := s.jobClient.IsJobOpen(ctx, req.JobID)
	if err != nil {
		return nil, errors.New("failed to verify job status")
	}
	if !isOpen {
		return nil, errors.New("job is not open for applications")
	}

	application := &domain.JobApplication{
		ID:          uuid.New(),
		JobID:       req.JobID,
		CandidateID: candidateID,
		Status:      "applied",
		CoverLetter: utils.SanitizeString(req.CoverLetter),
		AppliedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.applicationRepo.Create(ctx, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (s *CandidateService) GetApplications(ctx context.Context, candidateID uuid.UUID, userID uuid.UUID) ([]domain.JobApplication, error) {
	candidate, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil {
		return nil, err
	}

	if candidate.UserID != userID {
		return nil, errors.New("you can only view your own applications")
	}

	return s.applicationRepo.GetByCandidateID(ctx, candidateID)
}

func (s *CandidateService) ValidateUserPermissions(ctx context.Context, token string, requiredRole string) (*domain.UserInfo, error) {
	userInfo, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if requiredRole != "" && userInfo.Role != requiredRole {
		return nil, errors.New("insufficient permissions")
	}

	return userInfo, nil
}

func (s *CandidateService) loadCandidateRelations(ctx context.Context, candidate *domain.Candidate) error {
	skills, err := s.candidateSkillRepo.GetByCandidateID(ctx, candidate.ID)
	if err != nil {
		return err
	}
	candidate.Skills = skills

	workExps, err := s.workExpRepo.GetByCandidateID(ctx, candidate.ID)
	if err != nil {
		return err
	}
	candidate.WorkExperiences = workExps

	education, err := s.educationRepo.GetByCandidateID(ctx, candidate.ID)
	if err != nil {
		return err
	}
	candidate.Education = education

	resumes, err := s.resumeRepo.GetByCandidateID(ctx, candidate.ID)
	if err != nil {
		return err
	}
	candidate.Resumes = resumes

	applications, err := s.applicationRepo.GetByCandidateID(ctx, candidate.ID)
	if err != nil {
		return err
	}
	candidate.Applications = applications

	return nil
}

func (s *CandidateService) processResumeWithAI(ctx context.Context, resume *domain.Resume) {
	extractedText, err := s.aiService.ExtractTextFromResume(ctx, resume.FilePath)
	if err != nil {
		return
	}

	resume.ExtractedText = extractedText
	resume.AIProcessed = true
	resume.UpdatedAt = time.Now()

	s.resumeRepo.Update(ctx, resume)

	processedData, err := s.aiService.ProcessResumeData(ctx, extractedText)
	if err != nil {
		return
	}

	s.autoFillCandidateData(ctx, resume.CandidateID, processedData)
}

func (s *CandidateService) autoFillCandidateData(ctx context.Context, candidateID uuid.UUID, data *domain.ProcessedResumeData) {
	// Esta função poderia automaticamente preencher dados do candidato baseado na análise de IA
	// Por simplicidade, não implementaremos toda a lógica aqui
}
