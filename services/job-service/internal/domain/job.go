package domain

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text;not null"`
	Requirements string    `json:"requirements" gorm:"type:text"`
	Location    string     `json:"location"`
	SalaryMin   *float64   `json:"salary_min" gorm:"type:decimal(10,2)"`
	SalaryMax   *float64   `json:"salary_max" gorm:"type:decimal(10,2)"`
	Status      string     `json:"status" gorm:"not null;default:'open'"`
	CreatedBy   uuid.UUID  `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Skills      []JobSkill `json:"skills,omitempty" gorm:"foreignKey:JobID"`
}

type JobSkill struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	JobID         uuid.UUID `json:"job_id" gorm:"type:uuid;not null"`
	SkillID       uuid.UUID `json:"skill_id" gorm:"type:uuid;not null"`
	RequiredLevel string    `json:"required_level"`
	IsRequired    bool      `json:"is_required" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	Skill         *Skill    `json:"skill,omitempty" gorm:"foreignKey:SkillID"`
}

type Skill struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type JobStatus string

const (
	JobStatusOpen   JobStatus = "open"
	JobStatusClosed JobStatus = "closed"
)

func (j *Job) TableName() string {
	return "jobs"
}

func (js *JobSkill) TableName() string {
	return "job_skills"
}

func (s *Skill) TableName() string {
	return "skills"
}

func (j *Job) IsOpen() bool {
	return j.Status == string(JobStatusOpen)
}

func (j *Job) IsClosed() bool {
	return j.Status == string(JobStatusClosed)
}

type CreateJobRequest struct {
	Title        string              `json:"title" binding:"required"`
	Description  string              `json:"description" binding:"required"`
	Requirements string              `json:"requirements"`
	Location     string              `json:"location"`
	SalaryMin    *float64            `json:"salary_min"`
	SalaryMax    *float64            `json:"salary_max"`
	Skills       []CreateJobSkillRequest `json:"skills"`
}

type CreateJobSkillRequest struct {
	SkillID       uuid.UUID `json:"skill_id" binding:"required"`
	RequiredLevel string    `json:"required_level" binding:"required,oneof=beginner intermediate advanced expert"`
	IsRequired    bool      `json:"is_required"`
}

type UpdateJobRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Requirements string   `json:"requirements"`
	Location     string   `json:"location"`
	SalaryMin    *float64 `json:"salary_min"`
	SalaryMax    *float64 `json:"salary_max"`
}

type UpdateJobStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=open closed"`
}

type JobListFilter struct {
	Status   string
	Location string
	Title    string
	MinSalary *float64
	MaxSalary *float64
}

type JobResponse struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Requirements string           `json:"requirements"`
	Location    string            `json:"location"`
	SalaryMin   *float64          `json:"salary_min"`
	SalaryMax   *float64          `json:"salary_max"`
	Status      string            `json:"status"`
	CreatedBy   uuid.UUID         `json:"created_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Skills      []JobSkillResponse `json:"skills,omitempty"`
}

type JobSkillResponse struct {
	Skill         SkillResponse `json:"skill"`
	RequiredLevel string        `json:"required_level"`
	IsRequired    bool          `json:"is_required"`
}

type SkillResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
}
