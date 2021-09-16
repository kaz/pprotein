package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kaz/pprotein/integration/echov4"
	"github.com/labstack/echo/v4"
)

var (
	db *sql.DB
)

func handle(c echo.Context) error {
	time.Sleep(time.Duration(rand.Int63n(4096)) * time.Millisecond)

	if _, err := db.Exec("INSERT INTO mock.mock VALUES (?, NOW())", c.Param("id")); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := db.QueryRow("SELECT * FROM mock.mock WHERE id = ?", c.Param("id")).Err(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.String(http.StatusOK, "ok")
}

func start() error {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS mock"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS mock.mock (id VARCHAR(250) PRIMARY KEY, created DATETIME)"); err != nil {
		return err
	}

	e := echo.New()
	echov4.Integrate(e)

	e.Any("/api/mock/:id", handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "19999"
	}
	return e.Start(":" + port)
}

func startRequest() {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				if _, err := http.Get(fmt.Sprintf("http://%s/api/mock/%s", os.Getenv("REQUEST_HOST"), strconv.FormatInt(time.Now().UnixNano(), 36))); err != nil {
					log.Printf("request failed: %v", err)
				}
				time.Sleep(time.Duration(rand.Int63n(4096)) * time.Millisecond)
			}
		}()
	}
}

func main() {
	startRequest()

	if err := start(); err != nil {
		panic(err)
	}
}
