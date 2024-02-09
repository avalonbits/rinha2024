package main

import (
	"log"
	"net/http"
	"os"

	"github.com/avalonbits/rinha2024/endpoints/api"
	"github.com/avalonbits/rinha2024/storage/datastore"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := datastore.GetDB(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
