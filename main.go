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
		if len(os.Getenv("URL")) > 0 {
			fmt.Println("URL:", os.Getenv("URL"))
		}
		db, pgConnString, err := initStore()
		if err != nil {
			fmt.Println("err:", err)
			return c.HTML(http.StatusOK, pgConnString+":"+err.Error())
		}
		m := make([]map[string]interface{}, 0)
		if err = db.Table("user").Where("id= ?", 1).Find(&m); err != nil {
			fmt.Println("Find err:", err)
		}
		fmt.Println("user:", m)
		fmt.Println("我是admin22222222222222222")
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
	/*pgConnString = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", "124.223.101.122",
	"5432", "test", "root", "sl6387570506")*/
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
