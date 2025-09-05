package infrastructure

import (
	"context"
	"strings"

	"recruitment-system/services/candidate-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkExperienceRepositoryImpl struct {
	db *gorm.DB
}

func NewWorkExperienceRepository(db *gorm.DB) domain.WorkExperienceRepository {
	return &WorkExperienceRepositoryImpl{db: db}
}

func (r *WorkExperienceRepositoryImpl) Create(ctx context.Context, workExperience *domain.WorkExperience) error {
	return r.db.WithContext(ctx).Create(workExperience).Error
}

func (r *WorkExperienceRepositoryImpl) GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]domain.WorkExperience, error) {
	var workExperiences []domain.WorkExperience
	err := r.db.WithContext(ctx).
		Where("candidate_id = ?", candidateID).
		Order("start_date DESC").
		Find(&workExperiences).Error
	return workExperiences, err
}

func (r *WorkExperienceRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.WorkExperience, error) {
	var workExperience domain.WorkExperience
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&workExperience).Error
	if err != nil {
		return nil, err
	}
	return &workExperience, nil
}

func (r *WorkExperienceRepositoryImpl) Update(ctx context.Context, workExperience *domain.WorkExperience) error {
	return r.db.WithContext(ctx).Save(workExperience).Error
}

func (r *WorkExperienceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.WorkExperience{}, id).Error
}

type EducationRepositoryImpl struct {
	db *gorm.DB
}

func NewEducationRepository(db *gorm.DB) domain.EducationRepository {
	return &EducationRepositoryImpl{db: db}
}

func (r *EducationRepositoryImpl) Create(ctx context.Context, education *domain.Education) error {
	return r.db.WithContext(ctx).Create(education).Error
}

func (r *EducationRepositoryImpl) GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]domain.Education, error) {
	var education []domain.Education
	err := r.db.WithContext(ctx).
		Where("candidate_id = ?", candidateID).
		Order("start_date DESC").
		Find(&education).Error
	return education, err
}

func (r *EducationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Education, error) {
	var education domain.Education
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&education).Error
	if err != nil {
		return nil, err
	}
	return &education, nil
}

func (r *EducationRepositoryImpl) Update(ctx context.Context, education *domain.Education) error {
	return r.db.WithContext(ctx).Save(education).Error
}

func (r *EducationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Education{}, id).Error
}

type ResumeRepositoryImpl struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) domain.ResumeRepository {
	return &ResumeRepositoryImpl{db: db}
}

func (r *ResumeRepositoryImpl) Create(ctx context.Context, resume *domain.Resume) error {
	return r.db.WithContext(ctx).Create(resume).Error
}

func (r *ResumeRepositoryImpl) GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]domain.Resume, error) {
	var resumes []domain.Resume
	err := r.db.WithContext(ctx).
		Where("candidate_id = ?", candidateID).
		Order("created_at DESC").
		Find(&resumes).Error
	return resumes, err
}

func (r *ResumeRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Resume, error) {
	var resume domain.Resume
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepositoryImpl) Update(ctx context.Context, resume *domain.Resume) error {
	return r.db.WithContext(ctx).Save(resume).Error
}

func (r *ResumeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Resume{}, id).Error
}

type JobApplicationRepositoryImpl struct {
	db *gorm.DB
}

func NewJobApplicationRepository(db *gorm.DB) domain.JobApplicationRepository {
	return &JobApplicationRepositoryImpl{db: db}
}

func (r *JobApplicationRepositoryImpl) Create(ctx context.Context, application *domain.JobApplication) error {
	return r.db.WithContext(ctx).Create(application).Error
}

func (r *JobApplicationRepositoryImpl) GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]domain.JobApplication, error) {
	var applications []domain.JobApplication
	err := r.db.WithContext(ctx).
		Preload("Job").
		Where("candidate_id = ?", candidateID).
		Order("applied_at DESC").
		Find(&applications).Error
	return applications, err
}

func (r *JobApplicationRepositoryImpl) GetByJobID(ctx context.Context, jobID uuid.UUID) ([]domain.JobApplication, error) {
	var applications []domain.JobApplication
	err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("applied_at DESC").
		Find(&applications).Error
	return applications, err
}

func (r *JobApplicationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.JobApplication, error) {
	var application domain.JobApplication
	err := r.db.WithContext(ctx).
		Preload("Job").
		Where("id = ?", id).
		First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *JobApplicationRepositoryImpl) Update(ctx context.Context, application *domain.JobApplication) error {
	return r.db.WithContext(ctx).Save(application).Error
}

func (r *JobApplicationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.JobApplication{}, id).Error
}

func (r *JobApplicationRepositoryImpl) ExistsByCandidateAndJob(ctx context.Context, candidateID, jobID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.JobApplication{}).
		Where("candidate_id = ? AND job_id = ?", candidateID, jobID).
		Count(&count).Error
	return count > 0, err
}

type SkillRepositoryImpl struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) domain.SkillRepository {
	return &SkillRepositoryImpl{db: db}
}

func (r *SkillRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Skill, error) {
	var skill domain.Skill
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&skill).Error
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *SkillRepositoryImpl) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Skill, error) {
	var skills []domain.Skill
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&skills).Error
	return skills, err
}

func (r *SkillRepositoryImpl) List(ctx context.Context, category, search string, offset, limit int) ([]*domain.Skill, int64, error) {
	var skills []*domain.Skill
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Skill{})

	if category != "" {
		query = query.Where("LOWER(category) = ?", strings.ToLower(category))
	}

	if search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&skills).Error
	return skills, total, err
}
