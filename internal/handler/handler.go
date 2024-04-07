package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
}

type Deps struct{}

func NewHandler(deps Deps) *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())

	srv := NewStrictHandler(h, nil)
	RegisterHandlers(e, srv)

	return e
}
