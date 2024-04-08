package domain

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `database:"id"`
	Surname  string    `database:"surname"`
	Name     string    `database:"name"`
	Patronyc string    `database:"patronyc"`
	Email    string    `database:"email"`
	Password string    `database:"hash_password"`
	Verified bool      `database:"verified"`
}

type UserUpdate struct {
	Surname  *string
	Name     *string
	Patronyc *string
	Email    *string
	Password *string
	Verified *bool
}

func (u *UserUpdate) Validate() bool {
	if u.Surname == nil && u.Name == nil && u.Patronyc == nil && u.Email == nil &&
		u.Password == nil && u.Verified == nil {
		return false
	}
	return true
}
