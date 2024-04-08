package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/repository"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/hasher"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/tokens"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	usersRepo    repository.Users
	rdb          *redis.Client
	tokenManager tokens.TokenManagerInterface
	hasher       hasher.HasherInterface
}

func NewAuthService(usersRepo repository.Users, rdb *redis.Client,
	tokenManager tokens.TokenManagerInterface, hasher hasher.HasherInterface) *AuthService {
	return &AuthService{
		usersRepo:    usersRepo,
		rdb:          rdb,
		tokenManager: tokenManager,
		hasher:       hasher,
	}
}

type SendEmailVerificationMessageTask struct {
	Email string `json:"email"`
}

func (s *AuthService) writeTaskSendEmailVerificationMessage(ctx context.Context, email string) error {
	data, err := json.Marshal(SendEmailVerificationMessageTask{
		Email: email,
	})
	if err != nil {
		logrus.Errorf("error marshaling json send email verification task: %s", err)
		return ErrInternal
	}
	_, err = s.rdb.RPush(ctx, "queue:verification:email", data).Result()
	if err != nil {
		logrus.Errorf("error write task send email verification message into broker: %s", err)
	}
	return nil
}

func (s *AuthService) SendEmailVerificationMessage(ctx context.Context, id uuid.UUID,
	email string) error {
	user, err := s.usersRepo.Get(ctx, id)
	if err != nil {
		logrus.Errorf("error getting user from users repo when sending email verification message: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return ErrUserNotFound
		}
		return ErrInternal
	}
	if user.Verified {
		return ErrUserNotFound
	}
	return s.writeTaskSendEmailVerificationMessage(ctx, email)
}

func (s *AuthService) SignUp(ctx context.Context, user domain.User) (uuid.UUID, error) {
	var id uuid.UUID
	_, err := s.usersRepo.GetByEmail(ctx, user.Email)
	if err == nil {
		logrus.Errorf("email already in use: %s", err)
		return id, ErrEmailAlreadyInUse
	}
	if err != nil && !errors.Is(repository.ErrUserNotFound, err) {
		logrus.Errorf("error getting user from users repo when signup: %s", err)
		return id, ErrInternal
	}

	user.Password = s.hasher.Hash(user.Password)
	id, err = s.usersRepo.Create(ctx, user)
	if err != nil {
		logrus.Errorf("error creating user into repository when signing up: %s", err)
		return id, ErrInternal
	}

	_ = s.writeTaskSendEmailVerificationMessage(ctx, user.Email)

	return id, nil
}

func (s *AuthService) SignIn(ctx context.Context, email string, password string) (string, error) {
	user, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		logrus.Errorf("error getting user from repo by email when signing in: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return "", nil
		}
		return "", ErrInternal
	}

	validPassword := s.hasher.Check(password, user.Password)
	if !validPassword {
		logrus.Errorf("invalid password when signing in")
		return "", ErrInvalidEmailOrPassword
	}

	accessToken, err := s.tokenManager.CreateAccessToken(user.Id)
	if err != nil {
		logrus.Errorf("error creating access token when signing in: %s", err)
		return "", ErrInternal
	}

	return accessToken, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	email, err := s.tokenManager.ParseEmailToken(token)
	if err != nil {
		logrus.Errorf("error parsing email token when verifying email: %s", err)
		if errors.Is(tokens.ErrTokenExpired, err) {
			return ErrTokenExpired
		}
		if errors.Is(tokens.ErrTokenInvalid, err) {
			return ErrTokenInvalid
		}
		return ErrInternal
	}

	user, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		logrus.Errorf("error getting user from repo when verifying email: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return ErrUserNotFound
		}
		return ErrInternal
	}

	if user.Verified {
		return nil
	}

	verified := true
	data := domain.UserUpdate{
		Verified: &verified,
	}
	user, err = s.usersRepo.Update(ctx, user.Id, data)
	if err != nil {
		logrus.Errorf("error updating user into repo when verifying email: %s", err)
		return ErrInternal
	}

	return nil
}
