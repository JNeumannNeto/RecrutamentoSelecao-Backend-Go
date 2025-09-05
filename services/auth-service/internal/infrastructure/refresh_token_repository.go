package infrastructure

import (
	"context"
	"time"

	"recruitment-system/services/auth-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt int64     `gorm:"not null"`
	CreatedAt time.Time
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

type RefreshTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) domain.RefreshTokenRepository {
	db.AutoMigrate(&RefreshToken{})
	return &RefreshTokenRepositoryImpl{db: db}
}

func (r *RefreshTokenRepositoryImpl) Store(ctx context.Context, userID uuid.UUID, token string, expiresAt int64) error {
	refreshToken := &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(refreshToken).Error
}

func (r *RefreshTokenRepositoryImpl) Get(ctx context.Context, token string) (uuid.UUID, error) {
	var refreshToken RefreshToken
	err := r.db.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now().Unix()).First(&refreshToken).Error
	if err != nil {
		return uuid.Nil, err
	}
	return refreshToken.UserID, nil
}

func (r *RefreshTokenRepositoryImpl) Delete(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&RefreshToken{}).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}

func (r *RefreshTokenRepositoryImpl) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at <= ?", time.Now().Unix()).Delete(&RefreshToken{}).Error
}
