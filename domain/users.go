package domain

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `database:"id"`
	Surname  string    `database:"surname"`
	Name     string    `database:"name"`
	Patronyc string    `database:"patronyc"`
	Email    string    `database:"email"`
	Password string    `databas:"hash_password"`
}
