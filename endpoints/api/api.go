package api

import (
	"net/http"

	"github.com/avalonbits/rinha2024/service/rinha"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *rinha.Service
}

func New(svc *rinha.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Transact(c echo.Context) error {
	return httpError(http.StatusInternalServerError, "Unimplemented")
}

func (h *Handler) AccountHistory(c echo.Context) error {
	return httpError(http.StatusInternalServerError, "Unimplemented")
}
func httpError(status int, msg string) *echo.HTTPError {
	return echo.NewHTTPError(status, msg)
}
