package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	RefreshToken string    `json:"refreshToken"`
	Id           uuid.UUID `json:"id"`
	ExpiredAt    time.Time `json:"expiredAt"`
}
