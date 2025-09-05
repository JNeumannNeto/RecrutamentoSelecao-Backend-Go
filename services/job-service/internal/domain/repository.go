package domain

import (
	"context"

	"github.com/google/uuid"
)

type JobRepository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter JobListFilter, offset, limit int) ([]*Job, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	GetByCreatedBy(ctx context.Context, createdBy uuid.UUID, offset, limit int) ([]*Job, int64, error)
}

type JobSkillRepository interface {
	CreateBatch(ctx context.Context, jobSkills []JobSkill) error
	GetByJobID(ctx context.Context, jobID uuid.UUID) ([]JobSkill, error)
	DeleteByJobID(ctx context.Context, jobID uuid.UUID) error
	Update(ctx context.Context, jobSkill *JobSkill) error
}

type SkillRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Skill, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]Skill, error)
	List(ctx context.Context, category, search string, offset, limit int) ([]*Skill, int64, error)
	Create(ctx context.Context, skill *Skill) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type AuthServiceClient interface {
	ValidateToken(ctx context.Context, token string) (*UserInfo, error)
}

type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
	Name  string    `json:"name"`
}
