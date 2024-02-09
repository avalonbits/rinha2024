package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
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
