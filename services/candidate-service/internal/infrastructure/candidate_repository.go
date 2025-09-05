package infrastructure

import (
	"context"

	"recruitment-system/services/candidate-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CandidateRepositoryImpl struct {
	db *gorm.DB
}

func NewCandidateRepository(db *gorm.DB) domain.CandidateRepository {
	return &CandidateRepositoryImpl{db: db}
}

func (r *CandidateRepositoryImpl) Create(ctx context.Context, candidate *domain.Candidate) error {
	return r.db.WithContext(ctx).Create(candidate).Error
}

func (r *CandidateRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Candidate, error) {
	var candidate domain.Candidate
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&candidate).Error
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (r *CandidateRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Candidate, error) {
	var candidate domain.Candidate
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		First(&candidate).Error
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (r *CandidateRepositoryImpl) Update(ctx context.Context, candidate *domain.Candidate) error {
	return r.db.WithContext(ctx).Save(candidate).Error
}

func (r *CandidateRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Candidate{}, id).Error
}

func (r *CandidateRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*domain.Candidate, int64, error) {
	var candidates []*domain.Candidate
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Candidate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&candidates).Error

	return candidates, total, err
}

func (r *CandidateRepositoryImpl) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Candidate{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}

type CandidateSkillRepositoryImpl struct {
	db *gorm.DB
}

func NewCandidateSkillRepository(db *gorm.DB) domain.CandidateSkillRepository {
	return &CandidateSkillRepositoryImpl{db: db}
}

func (r *CandidateSkillRepositoryImpl) Create(ctx context.Context, candidateSkill *domain.CandidateSkill) error {
	return r.db.WithContext(ctx).Create(candidateSkill).Error
}

func (r *CandidateSkillRepositoryImpl) GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]domain.CandidateSkill, error) {
	var candidateSkills []domain.CandidateSkill
	err := r.db.WithContext(ctx).
		Preload("Skill").
		Where("candidate_id = ?", candidateID).
		Find(&candidateSkills).Error
	return candidateSkills, err
}

func (r *CandidateSkillRepositoryImpl) Update(ctx context.Context, candidateSkill *domain.CandidateSkill) error {
	return r.db.WithContext(ctx).Save(candidateSkill).Error
}

func (r *CandidateSkillRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.CandidateSkill{}, id).Error
}

func (r *CandidateSkillRepositoryImpl) DeleteByCandidateIDAndSkillID(ctx context.Context, candidateID, skillID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("candidate_id = ? AND skill_id = ?", candidateID, skillID).
		Delete(&domain.CandidateSkill{}).Error
}
