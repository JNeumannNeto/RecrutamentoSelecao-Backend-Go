package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"recruitment-system/services/candidate-service/internal/domain"

	"github.com/google/uuid"
)

type AuthServiceClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthServiceClient(baseURL string) domain.AuthServiceClient {
	return &AuthServiceClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *AuthServiceClientImpl) ValidateToken(ctx context.Context, token string) (*domain.UserInfo, error) {
	url := fmt.Sprintf("%s/api/v1/auth/validate", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			UserID string `json:"user_id"`
			Email  string `json:"email"`
			Role   string `json:"role"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("token validation failed: %s", response.Error)
	}

	userID, err := uuid.Parse(response.Data.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	return &domain.UserInfo{
		ID:    userID,
		Email: response.Data.Email,
		Role:  response.Data.Role,
		Name:  response.Data.Email,
	}, nil
}

func (c *AuthServiceClientImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

type JobServiceClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewJobServiceClient(baseURL string) domain.JobServiceClient {
	return &JobServiceClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *JobServiceClientImpl) GetJobByID(ctx context.Context, jobID uuid.UUID) (*domain.JobInfo, error) {
	url := fmt.Sprintf("%s/api/v1/jobs/%s", c.baseURL, jobID.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Location    string `json:"location"`
			Status      string `json:"status"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("failed to get job: %s", response.Error)
	}

	jobUUID, err := uuid.Parse(response.Data.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid job ID format: %w", err)
	}

	return &domain.JobInfo{
		ID:          jobUUID,
		Title:       response.Data.Title,
		Description: response.Data.Description,
		Location:    response.Data.Location,
		Status:      response.Data.Status,
	}, nil
}

func (c *JobServiceClientImpl) IsJobOpen(ctx context.Context, jobID uuid.UUID) (bool, error) {
	job, err := c.GetJobByID(ctx, jobID)
	if err != nil {
		return false, err
	}
	return job.Status == "open", nil
}

type FileStorageServiceImpl struct {
	uploadDir string
}

func NewFileStorageService(uploadDir string) domain.FileStorageService {
	return &FileStorageServiceImpl{
		uploadDir: uploadDir,
	}
}

func (s *FileStorageServiceImpl) SaveFile(ctx context.Context, file *multipart.FileHeader, candidateID uuid.UUID) (string, error) {
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	candidateDir := filepath.Join(s.uploadDir, candidateID.String())
	if err := os.MkdirAll(candidateDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create candidate directory: %w", err)
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(candidateDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return filePath, nil
}

func (s *FileStorageServiceImpl) DeleteFile(ctx context.Context, filePath string) error {
	return os.Remove(filePath)
}

func (s *FileStorageServiceImpl) GetFileURL(ctx context.Context, filePath string) (string, error) {
	return filePath, nil
}

type AIServiceImpl struct {
	apiURL string
	apiKey string
}

func NewAIService(apiURL, apiKey string) domain.AIService {
	return &AIServiceImpl{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

func (s *AIServiceImpl) ExtractTextFromResume(ctx context.Context, filePath string) (string, error) {
	return "Extracted text from resume (mock implementation)", nil
}

func (s *AIServiceImpl) ProcessResumeData(ctx context.Context, extractedText string) (*domain.ProcessedResumeData, error) {
	return &domain.ProcessedResumeData{
		Skills: []string{"Go", "Python", "JavaScript"},
		WorkExperiences: []struct {
			CompanyName string `json:"company_name"`
			Position    string `json:"position"`
			Description string `json:"description"`
			StartDate   string `json:"start_date"`
			EndDate     string `json:"end_date"`
			IsCurrent   bool   `json:"is_current"`
		}{
			{
				CompanyName: "Tech Company",
				Position:    "Software Developer",
				Description: "Developed web applications",
				StartDate:   "2020-01-01",
				EndDate:     "2023-12-31",
				IsCurrent:   false,
			},
		},
		Education: []struct {
			Institution  string  `json:"institution"`
			Degree       string  `json:"degree"`
			FieldOfStudy string  `json:"field_of_study"`
			StartDate    string  `json:"start_date"`
			EndDate      string  `json:"end_date"`
			IsCurrent    bool    `json:"is_current"`
			GPA          float64 `json:"gpa"`
		}{
			{
				Institution:  "University",
				Degree:       "Bachelor",
				FieldOfStudy: "Computer Science",
				StartDate:    "2016-01-01",
				EndDate:      "2019-12-31",
				IsCurrent:    false,
				GPA:          3.8,
			},
		},
	}, nil
}
