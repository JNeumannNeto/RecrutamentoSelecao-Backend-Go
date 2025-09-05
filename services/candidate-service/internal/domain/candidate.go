package domain

import (
	"time"

	"github.com/google/uuid"
)

type Candidate struct {
	ID              uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID         `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	Phone           string            `json:"phone"`
	Address         string            `json:"address"`
	DateOfBirth     *time.Time        `json:"date_of_birth" gorm:"type:date"`
	LinkedinURL     string            `json:"linkedin_url"`
	GithubURL       string            `json:"github_url"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	User            *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Skills          []CandidateSkill  `json:"skills,omitempty" gorm:"foreignKey:CandidateID"`
	WorkExperiences []WorkExperience  `json:"work_experiences,omitempty" gorm:"foreignKey:CandidateID"`
	Education       []Education       `json:"education,omitempty" gorm:"foreignKey:CandidateID"`
	Resumes         []Resume          `json:"resumes,omitempty" gorm:"foreignKey:CandidateID"`
	Applications    []JobApplication  `json:"applications,omitempty" gorm:"foreignKey:CandidateID"`
}

type User struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
}

type CandidateSkill struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CandidateID       uuid.UUID `json:"candidate_id" gorm:"type:uuid;not null"`
	SkillID           uuid.UUID `json:"skill_id" gorm:"type:uuid;not null"`
	ProficiencyLevel  string    `json:"proficiency_level"`
	YearsOfExperience int       `json:"years_of_experience" gorm:"default:0"`
	CreatedAt         time.Time `json:"created_at"`
	Skill             *Skill    `json:"skill,omitempty" gorm:"foreignKey:SkillID"`
}

type Skill struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
}

type WorkExperience struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CandidateID uuid.UUID  `json:"candidate_id" gorm:"type:uuid;not null"`
	CompanyName string     `json:"company_name" gorm:"not null"`
	Position    string     `json:"position" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text"`
	StartDate   time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate     *time.Time `json:"end_date" gorm:"type:date"`
	IsCurrent   bool       `json:"is_current" gorm:"default:false"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Education struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CandidateID   uuid.UUID  `json:"candidate_id" gorm:"type:uuid;not null"`
	Institution   string     `json:"institution" gorm:"not null"`
	Degree        string     `json:"degree" gorm:"not null"`
	FieldOfStudy  string     `json:"field_of_study"`
	StartDate     time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate       *time.Time `json:"end_date" gorm:"type:date"`
	IsCurrent     bool       `json:"is_current" gorm:"default:false"`
	GPA           *float64   `json:"gpa" gorm:"type:decimal(3,2)"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Resume struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CandidateID   uuid.UUID `json:"candidate_id" gorm:"type:uuid;not null"`
	Filename      string    `json:"filename" gorm:"not null"`
	FilePath      string    `json:"file_path" gorm:"not null"`
	FileSize      int64     `json:"file_size" gorm:"not null"`
	MimeType      string    `json:"mime_type" gorm:"not null"`
	ExtractedText string    `json:"extracted_text,omitempty" gorm:"type:text"`
	AIProcessed   bool      `json:"ai_processed" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type JobApplication struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	JobID       uuid.UUID `json:"job_id" gorm:"type:uuid;not null"`
	CandidateID uuid.UUID `json:"candidate_id" gorm:"type:uuid;not null"`
	Status      string    `json:"status" gorm:"not null;default:'applied'"`
	CoverLetter string    `json:"cover_letter" gorm:"type:text"`
	AppliedAt   time.Time `json:"applied_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Job         *Job      `json:"job,omitempty" gorm:"foreignKey:JobID"`
}

type Job struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
}

func (c *Candidate) TableName() string {
	return "candidates"
}

func (u *User) TableName() string {
	return "users"
}

func (cs *CandidateSkill) TableName() string {
	return "candidate_skills"
}

func (s *Skill) TableName() string {
	return "skills"
}

func (we *WorkExperience) TableName() string {
	return "work_experiences"
}

func (e *Education) TableName() string {
	return "education"
}

func (r *Resume) TableName() string {
	return "resumes"
}

func (ja *JobApplication) TableName() string {
	return "job_applications"
}

func (j *Job) TableName() string {
	return "jobs"
}

type CreateCandidateRequest struct {
	Phone       string     `json:"phone"`
	Address     string     `json:"address"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	LinkedinURL string     `json:"linkedin_url"`
	GithubURL   string     `json:"github_url"`
}

type UpdateCandidateRequest struct {
	Phone       string     `json:"phone"`
	Address     string     `json:"address"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	LinkedinURL string     `json:"linkedin_url"`
	GithubURL   string     `json:"github_url"`
}

type AddSkillRequest struct {
	SkillID           uuid.UUID `json:"skill_id" binding:"required"`
	ProficiencyLevel  string    `json:"proficiency_level" binding:"required,oneof=beginner intermediate advanced expert"`
	YearsOfExperience int       `json:"years_of_experience"`
}

type AddWorkExperienceRequest struct {
	CompanyName string     `json:"company_name" binding:"required"`
	Position    string     `json:"position" binding:"required"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	IsCurrent   bool       `json:"is_current"`
}

type AddEducationRequest struct {
	Institution  string     `json:"institution" binding:"required"`
	Degree       string     `json:"degree" binding:"required"`
	FieldOfStudy string     `json:"field_of_study"`
	StartDate    time.Time  `json:"start_date" binding:"required"`
	EndDate      *time.Time `json:"end_date"`
	IsCurrent    bool       `json:"is_current"`
	GPA          *float64   `json:"gpa"`
}

type CreateJobApplicationRequest struct {
	JobID       uuid.UUID `json:"job_id" binding:"required"`
	CoverLetter string    `json:"cover_letter"`
}

type CandidateResponse struct {
	ID              uuid.UUID                   `json:"id"`
	User            UserResponse                `json:"user"`
	Phone           string                      `json:"phone"`
	Address         string                      `json:"address"`
	DateOfBirth     *time.Time                  `json:"date_of_birth"`
	LinkedinURL     string                      `json:"linkedin_url"`
	GithubURL       string                      `json:"github_url"`
	Skills          []CandidateSkillResponse    `json:"skills,omitempty"`
	WorkExperiences []WorkExperienceResponse    `json:"work_experiences,omitempty"`
	Education       []EducationResponse         `json:"education,omitempty"`
	Resumes         []ResumeResponse            `json:"resumes,omitempty"`
	Applications    []JobApplicationResponse    `json:"applications,omitempty"`
	CreatedAt       time.Time                   `json:"created_at"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
}

type CandidateSkillResponse struct {
	ID                uuid.UUID     `json:"id"`
	Skill             SkillResponse `json:"skill"`
	ProficiencyLevel  string        `json:"proficiency_level"`
	YearsOfExperience int           `json:"years_of_experience"`
}

type SkillResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
}

type WorkExperienceResponse struct {
	ID          uuid.UUID  `json:"id"`
	CompanyName string     `json:"company_name"`
	Position    string     `json:"position"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	IsCurrent   bool       `json:"is_current"`
}

type EducationResponse struct {
	ID           uuid.UUID  `json:"id"`
	Institution  string     `json:"institution"`
	Degree       string     `json:"degree"`
	FieldOfStudy string     `json:"field_of_study"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	IsCurrent    bool       `json:"is_current"`
	GPA          *float64   `json:"gpa"`
}

type ResumeResponse struct {
	ID          uuid.UUID `json:"id"`
	Filename    string    `json:"filename"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `json:"mime_type"`
	AIProcessed bool      `json:"ai_processed"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

type JobApplicationResponse struct {
	ID          uuid.UUID   `json:"id"`
	Job         JobResponse `json:"job"`
	Status      string      `json:"status"`
	CoverLetter string      `json:"cover_letter"`
	AppliedAt   time.Time   `json:"applied_at"`
}

type JobResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
}
