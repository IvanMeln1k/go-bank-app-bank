package handler

import (
	"errors"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/sirupsen/logrus"
)

func (h *Handler) SignUp(ctx echo.Context) error {
	var userData UserWithPassword
	if err := ctx.Bind(&userData); err != nil {
		return httpBadRequest()
	}

	id, err := h.services.Auth.SignUp(ctx.Request().Context(), domain.User{
		Surname:  userData.Surname,
		Name:     userData.Name,
		Patronyc: userData.Patronyc,
		Email:    string(userData.Email),
		Password: userData.Password,
	})
	if err != nil {
		logrus.Errorf("error signup (handler): %s", err)
		if errors.Is(service.ErrEmailAlreadyInUse, err) {
			return ctx.JSON(409, Message{
				Message: "Email already in use",
			})
		}
		return httpInternalError()
	}

	return ctx.JSON(200, ReturnId{
		Id: id,
	})
}

func (h *Handler) SignIn(ctx echo.Context) error {
	var user AuthSchema
	if err := ctx.Bind(&user); err != nil {
		return httpBadRequest()
	}

	token, err := h.services.Auth.SignIn(ctx.Request().Context(), string(user.Email), user.Password)
	if err != nil {
		logrus.Errorf("error sign up (handler): %s", err)
		if errors.Is(service.ErrInvalidEmailOrPassword, err) {
			return ctx.JSON(401, Message{
				Message: "Invalid email or password",
			})
		}
		return httpInternalError()
	}

	return ctx.JSON(200, ReturnToken{
		Token: token,
	})
}

func (h *Handler) GetMe(ctx echo.Context) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}

	user, err := h.services.Auth.Get(ctx.Request().Context(), userId)
	if err != nil {
		logrus.Errorf("error getting user (handler): %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return httpErrUserNotFound()
		}
		return httpInternalError()
	}

	return ctx.JSON(200, map[string]interface{}{
		"User": User{
			Email:    types.Email(user.Email),
			Id:       userId,
			Surname:  user.Surname,
			Name:     user.Name,
			Patronyc: user.Patronyc,
			Verified: user.Verified,
		},
	})
}

func (h *Handler) VerifyEmail(ctx echo.Context, params VerifyEmailParams) error {
	err := h.services.VerifyEmail(ctx.Request().Context(), params.Token)
	if err != nil {
		if errors.Is(service.ErrTokenExpired, err) {
			return echo.NewHTTPError(401, Message{
				Message: "Token is expired",
			})
		}
		if errors.Is(service.ErrTokenInvalid, err) {
			return echo.NewHTTPError(401, Message{
				Message: "Token is invalid",
			})
		}
		return echo.NewHTTPError(500, Message{
			Message: "Internal server error",
		})
	}

	return ctx.JSON(200, Message{
		Message: "ok",
	})
}

func (h *Handler) ResendVerify(ctx echo.Context) error {
	userId, err := h.authorization(ctx)
	if err != nil {
		return err
	}

	err = h.services.Auth.SendEmailVerificationMessage(ctx.Request().Context(), userId)
	if err != nil {
		logrus.Errorf("error send email verification message when resending verify: %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return httpErrUserNotFound()
		}
		if errors.Is(service.ErrEmailAlreadyVerified, err) {
			return echo.NewHTTPError(409, Message{
				Message: "Email already verified",
			})
		}
		return httpInternalError()
	}

	return ctx.JSON(200, Message{
		Message: "ok",
	})
}
