package handler

import (
	"github.com/IvanMeln1k/go-bank-app-bank/internal/service"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/tokens"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	tokenManager tokens.TokenManagerInterface
	services     *service.Service
}

type Deps struct {
	TokenManager tokens.TokenManagerInterface
	Services     *service.Service
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		tokenManager: deps.TokenManager,
		services:     deps.Services,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	// e.Use(h.getUserIdentityMiddleware())

	RegisterHandlers(e, h)

	return e
}
