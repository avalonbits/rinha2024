package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/avalonbits/rinha2024/endpoints/api"
	"github.com/avalonbits/rinha2024/service/rinha"
	"github.com/avalonbits/rinha2024/storage/datastore"
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

func main() {
	db, err := datastore.GetDB(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	svc := rinha.New(db)
	handlers := api.New(svc)

	// Echo instance
	e := echo.New()
	e.JSONSerializer = easyJsonSerializer{}

	// Routes
	e.POST("/clientes/:id/transacoes", handlers.Transact)
	e.GET("/clients/:id/extrato", handlers.AccountHistory)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

type easyJsonSerializer struct {
}

func (_ easyJsonSerializer) Serialize(c echo.Context, data any, indent string) error {
	var buf []byte
	var err error

	ejs, ok := data.(easyjson.Marshaler)
	if ok {
		buf, err = easyjson.Marshal(ejs)
	} else {
		buf, err = json.Marshal(data)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err = io.Copy(c.Response(), bytes.NewBuffer(buf))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (_ easyJsonSerializer) Deserialize(c echo.Context, data any) error {
	js, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ejs, ok := data.(easyjson.Unmarshaler)
	if ok {
		err = easyjson.Unmarshal(js, ejs)
	} else {
		err = json.Unmarshal(js, data)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
