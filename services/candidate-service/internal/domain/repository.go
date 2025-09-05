package domain

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

type CandidateRepository interface {
	Create(ctx context.Context, candidate *Candidate) error
	GetByID(ctx context.Context, id uuid.UUID) (*Candidate, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Candidate, error)
	Update(ctx context.Context, candidate *Candidate) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Candidate, int64, error)
	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}

type CandidateSkillRepository interface {
	Create(ctx context.Context, candidateSkill *CandidateSkill) error
	GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]CandidateSkill, error)
	Update(ctx context.Context, candidateSkill *CandidateSkill) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByCandidateIDAndSkillID(ctx context.Context, candidateID, skillID uuid.UUID) error
}

type WorkExperienceRepository interface {
	Create(ctx context.Context, workExperience *WorkExperience) error
	GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]WorkExperience, error)
	GetByID(ctx context.Context, id uuid.UUID) (*WorkExperience, error)
	Update(ctx context.Context, workExperience *WorkExperience) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EducationRepository interface {
	Create(ctx context.Context, education *Education) error
	GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]Education, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Education, error)
	Update(ctx context.Context, education *Education) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ResumeRepository interface {
	Create(ctx context.Context, resume *Resume) error
	GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]Resume, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Resume, error)
	Update(ctx context.Context, resume *Resume) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type JobApplicationRepository interface {
	Create(ctx context.Context, application *JobApplication) error
	GetByCandidateID(ctx context.Context, candidateID uuid.UUID) ([]JobApplication, error)
	GetByJobID(ctx context.Context, jobID uuid.UUID) ([]JobApplication, error)
	GetByID(ctx context.Context, id uuid.UUID) (*JobApplication, error)
	Update(ctx context.Context, application *JobApplication) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByCandidateAndJob(ctx context.Context, candidateID, jobID uuid.UUID) (bool, error)
}

type SkillRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Skill, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]Skill, error)
	List(ctx context.Context, category, search string, offset, limit int) ([]*Skill, int64, error)
}

type FileStorageService interface {
	SaveFile(ctx context.Context, file *multipart.FileHeader, candidateID uuid.UUID) (string, error)
	DeleteFile(ctx context.Context, filePath string) error
	GetFileURL(ctx context.Context, filePath string) (string, error)
}

type AIService interface {
	ExtractTextFromResume(ctx context.Context, filePath string) (string, error)
	ProcessResumeData(ctx context.Context, extractedText string) (*ProcessedResumeData, error)
}

type ProcessedResumeData struct {
	Skills          []string `json:"skills"`
	WorkExperiences []struct {
		CompanyName string `json:"company_name"`
		Position    string `json:"position"`
		Description string `json:"description"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		IsCurrent   bool   `json:"is_current"`
	} `json:"work_experiences"`
	Education []struct {
		Institution  string  `json:"institution"`
		Degree       string  `json:"degree"`
		FieldOfStudy string  `json:"field_of_study"`
		StartDate    string  `json:"start_date"`
		EndDate      string  `json:"end_date"`
		IsCurrent    bool    `json:"is_current"`
		GPA          float64 `json:"gpa"`
	} `json:"education"`
}

type AuthServiceClient interface {
	ValidateToken(ctx context.Context, token string) (*UserInfo, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
}

type JobServiceClient interface {
	GetJobByID(ctx context.Context, jobID uuid.UUID) (*JobInfo, error)
	IsJobOpen(ctx context.Context, jobID uuid.UUID) (bool, error)
}

type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
	Name  string    `json:"name"`
}

type JobInfo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
}
