package domain

import (
	"time"

	"github.com/google/uuid"
)

// User — корневая сущность идентичности.
// user-info ссылается на неё по user_id (UUID).
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	IsBlocked    bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RefreshSession — долговременный сеанс (для refresh-токена).
// В БД храним только sha256(refresh) в TokenHash.
type RefreshSession struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash []byte
	UserAgent *string
	IP        *string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// PasswordReset — одноразовый сброс пароля по OTP.
type PasswordReset struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	OTPHash   string
	ExpiresAt time.Time
	UsedAt    *time.Time
}
