package domain

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `db:"id"`
	Surname  string    `db:"surname"`
	Name     string    `db:"name"`
	Patronyc string    `db:"patronyc"`
	Email    string    `db:"email"`
	Password string    `db:"hash_password"`
	Verified bool      `db:"verified"`
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
