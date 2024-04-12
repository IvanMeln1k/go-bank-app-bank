package handler

import "github.com/labstack/echo/v4"

func httpBadRequest() error {
	return echo.NewHTTPError(400, Message{
		Message: "Bad request",
	})
}

func httpInternalError() error {
	return echo.NewHTTPError(500, Message{
		Message: "Internal server error",
	})
}

func httpErrUserNotFound() error {
	return echo.NewHTTPError(404, Message{
		Message: "User not found",
	})
}

func httpErrAccountNotFound() error {
	return echo.NewHTTPError(404, Message{
		Message: "Account not found",
	})
}
