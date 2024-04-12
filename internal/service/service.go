package service

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/broker"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/repository"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/hasher"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/tokens"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/transactions"
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
	ErrInsufficientFunds      = errors.New("insufficient funds in the account")
	ErrTooManyAccounts        = errors.New("accounts can't be more 3")
	ErrEmailAlreadyVerified   = errors.New("email already verified")
)

type Auth interface {
	SignUp(ctx context.Context, user domain.User) (uuid.UUID, error)
	SignIn(ctx context.Context, email string, password string) (string, error)
	SendEmailVerificationMessage(ctx context.Context, id uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type Accounts interface {
	Create(ctx context.Context, userId uuid.UUID, account domain.Account) (uuid.UUID, error)
	Get(ctx context.Context, userId uuid.UUID, id uuid.UUID) (domain.Account, error)
	GetAll(ctx context.Context, userId uuid.UUID) ([]domain.Account, error)
	Delete(ctx context.Context, userId uuid.UUID, id uuid.UUID) error
	Transfer(ctx context.Context, userId uuid.UUID, id uuid.UUID, to uuid.UUID, amount int) error
}

type Machines interface {
	CashOut(ctx context.Context, id uuid.UUID, userId uuid.UUID, accountId uuid.UUID,
		amount int) error
	Deposit(ctx context.Context, id uuid.UUID, userId uuid.UUID, accountId uuid.UUID,
		amount int) error
}

type Service struct {
	Auth
	Accounts
	Machines
}

type Deps struct {
	Repos              *repository.Repository
	TokenManager       tokens.TokenManagerInterface
	RDB                *redis.Client
	Hasher             hasher.HasherInterface
	TransactionManager transactions.ManagerInterface
	Broker             broker.BrokerInterface
}

func NewService(deps Deps) *Service {
	return &Service{
		Auth: NewAuthService(deps.Repos.Users, deps.RDB, deps.TokenManager, deps.Hasher,
			deps.TransactionManager, deps.Broker),
		Accounts: NewAccountsService(deps.RDB, deps.Repos.Users, deps.Repos.Accounts,
			deps.TransactionManager),
		Machines: NewMachinesService(deps.Repos.Machines, deps.Repos.Accounts, deps.Repos.Users,
			deps.Broker),
	}
}
