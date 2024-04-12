package domain

import "github.com/google/uuid"

type Account struct {
	Id     uuid.UUID `db:"id"`
	Money  int       `db:"money"`
	UserId uuid.UUID `db:"user_id"`
}

type AccountUpdate struct {
	Money *int
}

func (a *AccountUpdate) Validate() bool {
	if a.Money == nil {
		return false
	}
	return true
}
