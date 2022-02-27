package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xormplus/xorm"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		db, pgConnString, err := initStore()
		if err != nil {
			return c.HTML(http.StatusOK, pgConnString+":"+err.Error())
		}
		defer db.Close()
		return c.HTML(http.StatusOK, pgConnString+"Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func initStore() (*xorm.Engine, string, error) {
	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	db, err := xorm.NewEngine("postgres", pgConnString)
	fmt.Println("connect address", pgConnString)
	if err != nil {
		logrus.Errorf("error creating db instance to %s: %s", pgConnString, err)
		return nil, pgConnString, err
	}
	db.SetMaxOpenConns(500)
	if err = db.Ping(); err != nil {
		logrus.Errorf("error creating db connection to %s: %s", pgConnString, err)
	}
	db.ShowSQL(true)
	return db, pgConnString, err
}
