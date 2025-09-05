package infrastructure

import (
	"context"
	"strings"

	"recruitment-system/services/job-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobRepositoryImpl struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) domain.JobRepository {
	return &JobRepositoryImpl{db: db}
}

func (r *JobRepositoryImpl) Create(ctx context.Context, job *domain.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *JobRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	var job domain.Job
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepositoryImpl) Update(ctx context.Context, job *domain.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *JobRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Job{}, id).Error
}

func (r *JobRepositoryImpl) List(ctx context.Context, filter domain.JobListFilter, offset, limit int) ([]*domain.Job, int64, error) {
	var jobs []*domain.Job
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Job{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Location != "" {
		query = query.Where("LOWER(location) LIKE ?", "%"+strings.ToLower(filter.Location)+"%")
	}

	if filter.Title != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(filter.Title)+"%")
	}

	if filter.MinSalary != nil {
		query = query.Where("salary_max >= ? OR salary_max IS NULL", *filter.MinSalary)
	}

	if filter.MaxSalary != nil {
		query = query.Where("salary_min <= ? OR salary_min IS NULL", *filter.MaxSalary)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&jobs).Error
	return jobs, total, err
}

func (r *JobRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Job{}).Where("id = ?", id).Update("status", status).Error
}

func (r *JobRepositoryImpl) GetByCreatedBy(ctx context.Context, createdBy uuid.UUID, offset, limit int) ([]*domain.Job, int64, error) {
	var jobs []*domain.Job
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Job{}).Where("created_by = ?", createdBy)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&jobs).Error
	return jobs, total, err
}
