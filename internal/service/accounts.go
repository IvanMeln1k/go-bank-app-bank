package service

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AccountsService struct {
	rdb          *redis.Client
	usersRepo    repository.Users
	accountsRepo repository.Accounts
}

func NewAccountsRepository(rdb *redis.Client, usersRepo repository.Users,
	accountsRepo repository.Accounts) *AccountsService {
	return &AccountsService{
		rdb:          rdb,
		usersRepo:    usersRepo,
		accountsRepo: accountsRepo,
	}
}

func (s *AccountsService) Create(ctx context.Context, userId uuid.UUID, account domain.Account) (uuid.UUID, error) {
	var id uuid.UUID

	user, err := s.usersRepo.Get(ctx, userId)
	if err != nil {
		logrus.Error("error getting user from repo when creating account: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return id, ErrUserNotFound
		}
		return id, ErrInternal
	}

	if !user.Verified {
		return id, ErrEmailNotVerified
	}

	id, err = s.accountsRepo.Create(ctx, userId, account)
	if err != nil {
		logrus.Errorf("error creating account into repo when creating account: %s", err)
		return id, ErrInternal
	}

	return id, nil
}

func (s *AccountsService) Get(ctx context.Context, userId uuid.UUID, id uuid.UUID) (uuid.Domain, error) {

}

func (s *AccountsService) Delete(ctx context.Context, userId uuid.UUID, id uuid.UUID) error {

}

func (s *AccountsService) Update(ctx context.Context, userId uuid.UUID, id uuid.UUID, data domain.AccountUpdate) (domain.Account, error) {

}

func (s *AccountsService) Transfer(ctx context.Context, userId uuid.UUID, id uuid.UUID, to uuid.UUID) error {

}
