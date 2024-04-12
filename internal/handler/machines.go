package handler

import (
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/service"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CashOut(ctx echo.Context, accountId openapi_types.UUID, params CashOutParams) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}

	var data CashOutJSONRequestBody
	if err := ctx.Bind(&data); err != nil {
		return httpBadRequest()
	}

	err = h.services.CashOut(ctx.Request().Context(), params.XMachineId, userId, accountId,
		int(data.Amount))
	if err != nil {
		logrus.Errorf("error cashout (handler): %s", err)
		if errors.Is(service.ErrMachineNotFound, err) {
			return echo.NewHTTPError(403, Message{
				Message: "Not enough rights",
			})
		}
		if errors.Is(service.ErrUserNotFound, err) {
			return httpErrUserNotFound()
		}
		if errors.Is(service.ErrAccountNotFound, err) {
			return httpErrAccountNotFound()
		}
		if errors.Is(service.ErrInsufficientFunds, err) {
			return echo.NewHTTPError(409, Message{
				Message: "Insufficient funds in the account",
			})
		}
		return httpInternalError()
	}

	return ctx.JSON(200, Message{
		Message: "ok",
	})
}

func (h *Handler) Deposit(ctx echo.Context, accountId openapi_types.UUID, params DepositParams) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}

	var data DepositJSONRequestBody
	if err := ctx.Bind(&data); err != nil {
		return httpBadRequest()
	}

	err = h.services.Deposit(ctx.Request().Context(), params.XMachineId, userId, accountId,
		int(data.Amount))
	if err != nil {
		logrus.Errorf("error cashout (handler): %s", err)
		if errors.Is(service.ErrMachineNotFound, err) {
			return echo.NewHTTPError(403, Message{
				Message: "Not enough rights",
			})
		}
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
