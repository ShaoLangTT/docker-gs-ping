package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"xorm.io/xorm"

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
			fmt.Println("err:", err)
			return c.HTML(http.StatusOK, pgConnString+":"+err.Error())
		}
		m := make([]map[string]interface{}, 0)
		db.Table("bg_user").Where("user_id= ?", "admin").Find(&m)
		fmt.Println("user:", m)
		fmt.Println("我是admin")
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
	fmt.Println("pgConnString:", pgConnString)
	/*	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", "124.223.101.122",
		"5432", "postgres", "postgres", "123456")*/
	engine, err := xorm.NewEngine("postgres", pgConnString)
	if err != nil {
		fmt.Println("NewEngine err:", err)
		return nil, pgConnString, err
	}
	setDB(engine)
	if err := engine.Ping(); err != nil {
		fmt.Println("Ping err:", err)
		return nil, pgConnString, err

	}
	return engine, pgConnString, err
}

func setDB(engine *xorm.Engine) {
	//engine.SetMaxIdleConns(3)
	//engine.SetMaxOpenConns(5)
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
}
