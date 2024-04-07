package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Users interface {
}

type Sessions interface {
}

type Accounts interface {
}

type Repository struct {
	Users
	Sessions
	Accounts
}

func NewRepository(pdb *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{}
}
