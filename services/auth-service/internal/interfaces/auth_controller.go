package interfaces

import (
	"net/http"

	"recruitment-system/services/auth-service/internal/application"
	"recruitment-system/services/auth-service/internal/domain"
	"recruitment-system/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	authService *application.AuthService
}

func NewAuthController(authService *application.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req domain.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	user, err := c.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Registration failed", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "User registered successfully", domain.UserInfo{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req domain.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	response, err := c.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Login failed", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Login successful", response)
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	response, err := c.authService.RefreshToken(ctx.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Token refresh failed", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Token refreshed successfully", response)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := c.authService.Logout(ctx.Request.Context(), userID); err != nil {
		utils.InternalServerErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Logout successful", nil)
}

func (c *AuthController) GetProfile(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := c.authService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		utils.NotFoundResponse(ctx, "User")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Profile retrieved successfully", domain.UserInfo{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	})
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req domain.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := c.authService.ChangePassword(ctx.Request.Context(), userID, req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Password change failed", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Password changed successfully", nil)
}

func (c *AuthController) ValidateToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Authorization header required", nil)
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := c.authService.ValidateToken(token)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Token is valid", gin.H{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
	})
}
