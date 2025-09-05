package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"recruitment-system/services/job-service/internal/domain"

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

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
		Role   string `json:"role"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
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

	var response ValidateTokenResponse
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

	userInfo := &domain.UserInfo{
		ID:    userID,
		Email: response.Data.Email,
		Role:  response.Data.Role,
		Name:  response.Data.Email,
	}

	return userInfo, nil
}
