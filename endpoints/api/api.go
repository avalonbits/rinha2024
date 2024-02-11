package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

//easyjson:json
type transactRequest struct {
	Value       int64  `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

func (r *transactRequest) validate(c echo.Context) error {
	if r.Value <= 0 {
		return fmt.Errorf("valor tem que ser positivo")
	}

	r.Type = strings.TrimSpace(r.Type)
	if r.Type == "" {
		return fmt.Errorf("tipo não definido")
	}
	if r.Type == "d" || r.Type == "D" {
		r.Value = -r.Value
	} else if r.Type != "c" && r.Type != "C" {
		return fmt.Errorf("tipo deve ser 'c' ou 'd'")
	}

	r.Description = strings.TrimSpace(r.Description)
	if r.Description == "" {
		return fmt.Errorf("descricao não definida")
	}
	if len(r.Description) > 10 {
		return fmt.Errorf("descricao muito longa")
	}

	return nil
}

func (h *Handler) Transact(c echo.Context) error {
	r := &transactRequest{}
	if err := h.validateRequest(c, r); err != nil {
		return httpError(http.StatusUnprocessableEntity, err.Error())
	}
	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpError(http.StatusUnprocessableEntity, "id inválido")
	}

	res, err := h.svc.Transact(c.Request().Context(), int64(cid), r.Value, r.Description)
	if err != nil {
		if errors.Is(err, rinha.OverLimitErr) {
			return httpError(http.StatusUnprocessableEntity, err.Error())
		}
		if errors.Is(err, rinha.NotFoundErr) {
			return httpError(http.StatusNotFound, err.Error())
		}
		return httpError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AccountHistory(c echo.Context) error {
	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpError(http.StatusBadRequest, "id inválido")
	}

	res, err := h.svc.AccountHistory(c.Request().Context(), int64(cid))
	if err != nil {
		if errors.Is(err, rinha.NotFoundErr) {
			return httpError(http.StatusNotFound, err.Error())
		}
		return httpError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func httpError(status int, msg string) *echo.HTTPError {
	return echo.NewHTTPError(status, msg)
}

type validator interface {
	validate(echo.Context) error
}

func (h *Handler) validateRequest(c echo.Context, req validator) error {
	var err error
	if err = c.Bind(req); err == nil {
		err = req.validate(c)
	}

	if err != nil {
		return httpError(http.StatusBadRequest, err.Error())
	}
	return nil
}
