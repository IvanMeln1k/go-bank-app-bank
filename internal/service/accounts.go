package service

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/repository"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/transactions"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AccountsService struct {
	rdb                *redis.Client
	usersRepo          repository.Users
	accountsRepo       repository.Accounts
	transactionManager transactions.ManagerInterface
}

func NewAccountsService(rdb *redis.Client, usersRepo repository.Users,
	accountsRepo repository.Accounts, transactionManager transactions.ManagerInterface) *AccountsService {
	return &AccountsService{
		rdb:                rdb,
		usersRepo:          usersRepo,
		accountsRepo:       accountsRepo,
		transactionManager: transactionManager,
	}
}

func (s *AccountsService) Create(ctx context.Context, userId uuid.UUID, account domain.Account) (uuid.UUID, error) {
	var id uuid.UUID

	user, err := s.usersRepo.Get(ctx, userId)
	if err != nil {
		logrus.Errorf("error getting user from repo when creating account: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return id, ErrUserNotFound
		}
		return id, ErrInternal
	}

	if !user.Verified {
		return id, ErrEmailNotVerified
	}

	accounts, err := s.accountsRepo.GetAll(ctx, userId)
	if err != nil {
		logrus.Errorf("error getting all accounts from repo where creating: %s", err)
		return uuid.UUID{}, ErrInternal
	}
	if len(accounts) >= 3 {
		logrus.Errorf("too many accounts for creating new: %s", err)
		return uuid.UUID{}, ErrTooManyAccounts
	}

	id, err = s.accountsRepo.Create(ctx, userId, account)
	if err != nil {
		logrus.Errorf("error creating account into repo when creating account: %s", err)
		return id, ErrInternal
	}

	return id, nil
}

func (s *AccountsService) GetAll(ctx context.Context, userId uuid.UUID) ([]domain.Account, error) {
	_, err := s.usersRepo.Get(ctx, userId)
	if err != nil {
		logrus.Errorf("error getting user from repo when get all accounts: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}

	accounts, err := s.accountsRepo.GetAll(ctx, userId)
	if err != nil {
		logrus.Errorf("error getting all accounts from repo when get all accounts: %s", err)
		return nil, ErrInternal
	}

	return accounts, nil
}

func (s *AccountsService) get(ctx context.Context, userId uuid.UUID, id uuid.UUID) (domain.Account, error) {
	account, err := s.accountsRepo.Get(ctx, id)
	if err != nil {
		logrus.Errorf("error getting account from repo when simple getting: %s", err)
		if errors.Is(repository.ErrAccountNotFound, err) {
			return account, ErrAccountNotFound
		}
		return account, ErrInternal
	}
	if account.UserId != userId {
		logrus.Errorf("error account %s doesn't belong user %s", id, userId)
		return account, ErrAccountNotFound
	}

	return account, nil
}

func (s *AccountsService) Get(ctx context.Context, userId uuid.UUID, id uuid.UUID) (domain.Account, error) {
	return s.get(ctx, userId, id)
}

func (s *AccountsService) Delete(ctx context.Context, userId uuid.UUID, id uuid.UUID) error {
	_, err := s.get(ctx, userId, id)
	if err != nil {
		return err
	}

	if err := s.accountsRepo.Delete(ctx, id); err != nil {
		logrus.Errorf("error deleting account from repo when deleting: %s", err)
		return ErrInternal
	}

	return nil
}

func (s *AccountsService) Transfer(ctx context.Context, userId uuid.UUID, id uuid.UUID,
	to uuid.UUID, amount int) error {
	account, err := s.get(ctx, userId, id)
	if err != nil {
		return err
	}

	if account.Money < amount {
		logrus.Errorf("insufficient funds in the account %s to transfer amount %d", id, account.Money)
		return ErrInsufficientFunds
	}

	accountTo, err := s.get(ctx, userId, to)
	if err != nil {
		logrus.Errorf("error get accountTo when transfering (service method): %s", err)
		if errors.Is(repository.ErrAccountNotFound, err) {
			return ErrAccountNotFound
		}
		return ErrInternal
	}

	err = s.transactionManager.Do(ctx, func(ctx context.Context) error {
		newMoneyFrom := account.Money - amount
		newMoneyTo := accountTo.Money + amount

		_, err := s.accountsRepo.Update(ctx, id, domain.AccountUpdate{
			Money: &newMoneyFrom,
		})
		if err != nil {
			return err
		}
		_, err = s.accountsRepo.Update(ctx, to, domain.AccountUpdate{
			Money: &newMoneyTo,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logrus.Errorf("error transfering transaction in service transfer method: %s", err)
		return ErrInternal
	}

	return nil
}
