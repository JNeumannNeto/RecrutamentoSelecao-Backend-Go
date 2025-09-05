package infrastructure

import (
	"context"
	"strings"

	"recruitment-system/services/job-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (r *SkillRepositoryImpl) Create(ctx context.Context, skill *domain.Skill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

func (r *SkillRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Skill{}).Where("LOWER(name) = ?", strings.ToLower(name)).Count(&count).Error
	return count > 0, err
}

type JobSkillRepositoryImpl struct {
	db *gorm.DB
}

func NewJobSkillRepository(db *gorm.DB) domain.JobSkillRepository {
	return &JobSkillRepositoryImpl{db: db}
}

func (r *JobSkillRepositoryImpl) CreateBatch(ctx context.Context, jobSkills []domain.JobSkill) error {
	if len(jobSkills) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&jobSkills).Error
}

func (r *JobSkillRepositoryImpl) GetByJobID(ctx context.Context, jobID uuid.UUID) ([]domain.JobSkill, error) {
	var jobSkills []domain.JobSkill
	err := r.db.WithContext(ctx).
		Preload("Skill").
		Where("job_id = ?", jobID).
		Find(&jobSkills).Error
	return jobSkills, err
}

func (r *JobSkillRepositoryImpl) DeleteByJobID(ctx context.Context, jobID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("job_id = ?", jobID).Delete(&domain.JobSkill{}).Error
}

func (r *JobSkillRepositoryImpl) Update(ctx context.Context, jobSkill *domain.JobSkill) error {
	return r.db.WithContext(ctx).Save(jobSkill).Error
}
