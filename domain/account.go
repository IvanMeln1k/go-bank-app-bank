package domain

import "github.com/google/uuid"

type Account struct {
	Id     uuid.UUID `database:"id"`
	Money  int       `database:"money"`
	UserId uuid.UUID `database:"user_id"`
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
