package main

import (
	"net/http"

	"github.com/avalonbits/rinha2024/endpoints/api"
	"github.com/labstack/echo/v4"
)

func main() {
	// Echo instance
	e := echo.New()
	handlers := api.New()

	// Routes
	e.POST("/clientes/:id/transacoes", handlers.Transact)
	e.GET("/clients/:id/extrato", handlers.AccountHistory)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
