package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/trace"

	_ "net/http/pprof"

	"github.com/avalonbits/rinha2024/endpoints/api"
	"github.com/avalonbits/rinha2024/service/rinha"
	"github.com/avalonbits/rinha2024/storage/datastore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mailru/easyjson"
)

var (
	traceF = flag.Bool("trace", false, "If true, writes tracing data for the server.")
)

func main() {
	flag.Parse()
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
	e.GET("/clientes/:id/extrato", handlers.AccountHistory, middleware.Gzip())

	// Start pprof server
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if *traceF {
		traceFile, err := os.Create("/tmp/trace.out")
		if err != nil {
			log.Fatal(err)
		}
		defer traceFile.Close()

		if err := trace.Start(traceFile); err != nil {
			log.Fatal(err)
		}
		defer trace.Stop()
	}

	go func() {
		// Start server
		e.Logger.Info(e.Start(":9999"))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	e.Shutdown(context.Background())
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
