package handler

import (
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/service"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateAccount(ctx echo.Context) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}

	accountId, err := h.services.Accounts.Create(ctx.Request().Context(), userId, domain.Account{
		Money:  0,
		UserId: userId,
	})
	if err != nil {
		logrus.Errorf("error creating account (handler): %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return httpErrUserNotFound()
		}
		if errors.Is(service.ErrEmailNotVerified, err) {
			return echo.NewHTTPError(403, "Email not verified")
		}
		if errors.Is(service.ErrTooManyAccounts, err) {
			return echo.NewHTTPError(409, "Too many accounts")
		}
		return httpInternalError()
	}

	return ctx.JSON(200, ReturnId{
		Id: accountId,
	})
}

func (h *Handler) GetAllAccounts(ctx echo.Context) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}

	accounts, err := h.services.Accounts.GetAll(ctx.Request().Context(), userId)
	if err != nil {
		logrus.Errorf("error get all accounts (handler): %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return httpErrUserNotFound()
		}
		return httpInternalError()
	}

	accountsReturn := make([]Account, len(accounts))
	for i, acc := range accounts {
		accountsReturn[i] = Account{
			Id:    acc.Id,
			Money: int32(acc.Money),
		}
	}

	return ctx.JSON(200, map[string]interface{}{
		"accounts": accountsReturn,
	})
}

func (h *Handler) GetAccountInfo(ctx echo.Context, accountId openapi_types.UUID) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}

	account, err := h.services.Accounts.Get(ctx.Request().Context(), userId, accountId)
	if err != nil {
		logrus.Errorf("error get account (handler): %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return httpErrUserNotFound()
		}
		if errors.Is(service.ErrAccountNotFound, err) {
			return httpErrAccountNotFound()
		}
		return httpInternalError()
	}

	return ctx.JSON(500, map[string]interface{}{
		"account": Account{
			Id:    accountId,
			Money: int32(account.Money),
		},
	})
}

func (h *Handler) DeleteAccount(ctx echo.Context, accountId openapi_types.UUID) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}

	err = h.services.Accounts.Delete(ctx.Request().Context(), userId, accountId)
	if err != nil {
		logrus.Errorf("error delete account (handler): %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return httpErrUserNotFound()
		}
		if errors.Is(service.ErrAccountNotFound, err) {
			return httpErrAccountNotFound()
		}
		return httpInternalError()
	}

	return ctx.JSON(200, Message{
		Message: "ok",
	})
}

func (h *Handler) Transfer(ctx echo.Context, accountId openapi_types.UUID) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}
	var transferInfo TransferInfo
	if err := ctx.Bind(&transferInfo); err != nil {
		return httpBadRequest()
	}

	err = h.services.Transfer(ctx.Request().Context(), userId, accountId, transferInfo.To,
		int(transferInfo.Amount))
	if err != nil {
		logrus.Errorf("error transfer (handler): %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return httpErrUserNotFound()
		}
		if errors.Is(service.ErrAccountNotFound, err) {
			return httpErrAccountNotFound()
		}
		if errors.Is(service.ErrInsufficientFunds, err) {
			return echo.NewHTTPError(409, "Insufficient funds in the account")
		}
		return httpInternalError()
	}

	return ctx.JSON(200, Message{
		Message: "ok",
	})
}
