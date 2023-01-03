package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aaunj/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth != "November 10, 2009" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

func setupRoute() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return next(c)
		}
	})

	e.GET("/health", healthHandler)
	g := e.Group("/expenses")
	g.POST("", AuthMiddleware(expense.CreateHandler))

	return e
}

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	expense.InitDB()

	r := setupRoute()

	go func() {
		if err := r.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			r.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Shutdown(ctx); err != nil {
		r.Logger.Fatal(err)
	}

}
