package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Репозиторий пользователей
type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (User, error)
	ByEmail(ctx context.Context, email string) (User, error)
	ByID(ctx context.Context, id uuid.UUID) (User, error)

	UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error
	SetLastLogin(ctx context.Context, id uuid.UUID, t time.Time) error
	SetBlocked(ctx context.Context, id uuid.UUID, blocked bool) error
}

// Репозиторий refresh-сессий
type RefreshRepository interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash []byte, ua, ip *string, exp time.Time) (RefreshSession, error)
	// Должен возвращать ТОЛЬКО активную сессию (revoked_at IS NULL && now < expires_at)
	ByHashActive(ctx context.Context, tokenHash []byte, now time.Time) (RefreshSession, error)

	RevokeByID(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}

// Репозиторий сброса пароля
type PasswordResetRepository interface {
	CreateOTP(ctx context.Context, userID uuid.UUID, otpHash string, exp time.Time) (PasswordReset, error)
	// Возвращает самый «свежий» валидный (used_at IS NULL && now < expires_at)
	FindValidByUser(ctx context.Context, userID uuid.UUID, now time.Time) (PasswordReset, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
}
