package handler

import (
	"errors"
	"strings"

	"github.com/IvanMeln1k/go-bank-app-bank/pkg/tokens"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) authorization(c echo.Context) (uuid.UUID, error) {
	authHeader := c.Request().Header.Get("Authorization")
	params := strings.Split(authHeader, " ")
	if len(params) != 2 || params[0] != "Bearer" {
		return uuid.UUID{}, echo.NewHTTPError(401, "No authorized")
	}
	userId, err := h.tokenManager.ParseAccessToken(params[1])
	if err != nil {
		logrus.Errorf("error parsing access token (handler): %s", err)
		if errors.Is(tokens.ErrTokenExpired, err) {
			return uuid.UUID{}, echo.NewHTTPError(401, Message{
				Message: "Token is expired",
			})
		}
		if errors.Is(tokens.ErrTokenInvalid, err) {
			return uuid.UUID{}, echo.NewHTTPError(401, Message{
				Message: "Token is invalid",
			})
		}
		return uuid.UUID{}, echo.NewHTTPError(500, Message{
			Message: "Internal server error",
		})
	}
	return userId, nil
}

// func (h *Handler) getUserIdentityMiddleware() echo.MiddlewareFunc {
// 	spec, _ := GetSwagger()
// 	m := middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
// 		Options: openapi3filter.Options{
// 			AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
// 				if ai.Scopes[0] == "user" || ai.Scopes[0] == "machine" {
// 					authHeader := ai.RequestValidationInput.Request.Header.Get("Authorization")
// 					params := strings.Split(authHeader, " ")
// 					if len(params) != 2 || params[0] != "Bearer" {
// 						return errors.New("no authorized")
// 					}
// 					userId, err := h.tokenManager.ParseAccessToken(params[1])
// 					if err != nil {
// 						if errors.Is(tokens.ErrTokenExpired, err) {
// 							return errors.New("token is expired")
// 						}
// 						if errors.Is(tokens.ErrTokenInvalid, err) {
// 							return errors.New("token is invalid")
// 						}
// 						return errors.New("internal server error")
// 					}
// 					eCxt := middleware.GetEchoContext(ctx)
// 					eCxt.Set("userId", userId)
// 				}
// 				return nil
// 			},
// 		},
// 	})
// 	return m
// }
