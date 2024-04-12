package service

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/broker"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MachinesService struct {
	machinesRepo repository.Machines
	accountsRepo repository.Accounts
	usersRepo    repository.Users
	broker       broker.BrokerInterface
}

func NewMachinesService(machinesRepo repository.Machines, accountsRepo repository.Accounts,
	usersRepo repository.Users, broker broker.BrokerInterface) *MachinesService {
	return &MachinesService{
		machinesRepo: machinesRepo,
		accountsRepo: accountsRepo,
		usersRepo:    usersRepo,
		broker:       broker,
	}
}

func (s *MachinesService) getAccount(ctx context.Context, id uuid.UUID,
	userId uuid.UUID) (domain.Account, error) {
	account, err := s.accountsRepo.Get(ctx, id)
	if err != nil {
		logrus.Errorf("error get account from repo in machine service: %s", err)
		if errors.Is(repository.ErrAccountNotFound, err) {
			return account, ErrAccountNotFound
		}
		return account, ErrInternal
	}
	if account.UserId != userId {
		logrus.Errorf("error account #%d doesn't belong user with id %d", id, userId)
		return account, ErrAccountNotFound
	}
	return account, nil
}

func (s *MachinesService) getUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.usersRepo.Get(ctx, id)
	if err != nil {
		logrus.Errorf("error getting user from repo in machine service: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return user, ErrUserNotFound
		}
		return user, ErrInternal
	}
	return user, nil
}

func (s *MachinesService) getMachine(ctx context.Context, id uuid.UUID) (domain.Machine, error) {
	machine, err := s.machinesRepo.Get(ctx, id)
	if err != nil {
		logrus.Errorf("error getting machine from repo: %s", err)
		if errors.Is(repository.ErrMachineNotFound, err) {
			return machine, ErrMachineNotFound
		}
		return machine, ErrInternal
	}
	return machine, nil
}

func (s *MachinesService) CashOut(ctx context.Context, id uuid.UUID, userId uuid.UUID,
	accountId uuid.UUID, amount int) error {
	_, err := s.getMachine(ctx, id)
	if err != nil {
		return err
	}

	user, err := s.getUser(ctx, userId)
	if err != nil {
		return err
	}

	account, err := s.getAccount(ctx, accountId, userId)
	if err != nil {
		return err
	}

	if account.Money < amount {
		logrus.Errorf("error insufficient funds in the account #%d for cash out", accountId)
		return ErrInsufficientFunds
	}

	newMoney := account.Money - amount
	_, err = s.accountsRepo.Update(ctx, accountId, domain.AccountUpdate{
		Money: &newMoney,
	})
	if err != nil {
		logrus.Errorf("error updating account when cashout: %s", err)
		return ErrInternal
	}

	s.broker.WriteCashoutTask(ctx, id, user.Email, accountId, amount, newMoney)

	return nil
}

func (s *MachinesService) Deposit(ctx context.Context, id uuid.UUID, userId uuid.UUID,
	accountId uuid.UUID, amount int) error {
	_, err := s.getMachine(ctx, id)
	if err != nil {
		return err
	}

	user, err := s.getUser(ctx, userId)
	if err != nil {
		return err
	}

	account, err := s.getAccount(ctx, accountId, userId)
	if err != nil {
		return err
	}

	newMoney := account.Money + amount
	_, err = s.accountsRepo.Update(ctx, accountId, domain.AccountUpdate{
		Money: &newMoney,
	})
	if err != nil {
		logrus.Errorf("error updating account when deposit: %s", err)
		return ErrInternal
	}

	s.broker.WriteDepositTask(ctx, id, user.Email, accountId, amount, newMoney)

	return nil
}
