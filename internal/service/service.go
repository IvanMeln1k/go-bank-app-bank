package service

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/repository"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/hasher"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/tokens"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInternal               = errors.New("error internal")
	ErrUserNotFound           = errors.New("user not found")
	ErrAccountNotFound        = errors.New("account not found")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrMachineNotFound        = errors.New("machine not found")
	ErrEmailAlreadyInUse      = errors.New("email already in use")
	ErrTokenExpired           = errors.New("token is expired")
	ErrTokenInvalid           = errors.New("token is invalid")
	ErrEmailNotVerified       = errors.New("email not verified")
)

type Auth interface {
	SignUp(ctx context.Context, user domain.User) (uuid.UUID, error)
	SignIn(ctx context.Context, email string, password string) (string, error)
	SendEmailVerificationMessage(ctx context.Context, id uuid.UUID, email string) error
	VerifyEmail(ctx context.Context, token string) error
}

type Users interface {
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type Accounts interface {
	Create(ctx context.Context, userId uuid.UUID, account domain.Account) (uuid.UUID, error)
	Get(ctx context.Context, userId uuid.UUID, id uuid.UUID) (uuid.Domain, error)
	Delete(ctx context.Context, userId uuid.UUID, id uuid.UUID) error
	Update(ctx context.Context, userId uuid.UUID, id uuid.UUID, data domain.AccountUpdate) (domain.Account, error)
	Transfer(ctx context.Context, userId uuid.UUID, id uuid.UUID, to uuid.UUID) error
}

type Machines interface {
	CashOut(ctx context.Context, id uuid.UUID, userId uuid.UUID, accountId uuid.UUID, amount int) error
	Deposit(ctx context.Context, id uuid.UUID, userId uuid.UUID, accountId uuid.UUID, amount int) error
}

type Service struct {
	Auth
	Users
	Accounts
	Machines
}

type Deps struct {
	Repos        *repository.Repository
	TokenManager tokens.TokenManagerInterface
	RDB          *redis.Client
	Hasher       hasher.HasherInterface
}

func NewService(deps Deps) *Service {
	return &Service{
		Auth:     NewAuthService(deps.Repos.Users, deps.RDB, deps.TokenManager, deps.Hasher),
		Accounts: NewAccountsRepository(deps.RDB, deps.Repos.Users, deps.Repos.Accounts),
	}
}
