package repository

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	usersTable    = "users"
	accountsTable = "accounts"
	machinesTable = "machines"
)

var (
	ErrInternal           = errors.New("internal error")
	ErrUserNotFound       = errors.New("user not found")
	ErrAccountNotFound    = errors.New("account not found")
	ErrSessionDoesntExist = errors.New("session doesn't exist")
	ErrMachineNotFound    = errors.New("machine not found")
)

type Users interface {
	Create(ctx context.Context, user domain.User) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, data domain.UserUpdate) (domain.User, error)
}

type Accounts interface {
	Create(ctx context.Context, userId uuid.UUID, account domain.Account) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Account, error)
	GetAll(ctx context.Context, userId uuid.UUID) ([]domain.Account, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, data domain.AccountUpdate) (domain.Account, error)
}

type Machines interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Machine, error)
}

type Repository struct {
	Users
	Accounts
	Machines
}

type Deps struct {
	DB *sqlx.DB
}

func NewRepository(deps Deps) *Repository {
	return &Repository{
		Users:    NewUsersRepository(deps.DB),
		Accounts: NewAccountsRepository(deps.DB),
		Machines: NewMachinesRepository(deps.DB),
	}
}
