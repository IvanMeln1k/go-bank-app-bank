package domain

import "github.com/google/uuid"

type Account struct {
	Id     uuid.UUID
	UserId uuid.UUID
	Money  int
}
