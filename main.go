package main

import (
	"github.com/google/uuid"
	"github.com/guseinovmg/imhioTest/handlers"
	"os/signal"
	"strconv"

	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

func setTokenAndCounter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			cookie = new(http.Cookie)
			cookie.Name = "token"
			uuidValue := uuid.New()
			cookie.Value = uuidValue.String()
			cookie.Expires = time.Now().Add(time.Minute)
			c.SetCookie(cookie)

			cookie = new(http.Cookie)
			cookie.Name = "counter"
			cookie.Value = "1"
			c.SetCookie(cookie)
		} else {
			cookie, err = c.Cookie("counter")
			if err != nil {
				return err
			}
			num, err := strconv.Atoi(cookie.Value)
			if err != nil {
				return err
			}
			cookie.Value = strconv.Itoa(num + 1)
			c.SetCookie(cookie)
		}
		return next(c)
	}
}

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	articleGroup := e.Group("/article")

	articleGroup.GET("/:id", handlers.GetArticleById)

	articleGroup.GET("/", handlers.GetArticlesByTag, setTokenAndCounter)

	articleGroup.POST("", handlers.CreateNewArticle)

	articleGroup.PUT("/:id", handlers.UpdateArticle)

	articleGroup.DELETE("/:id", handlers.DeleteArticle)

	e.Logger.Fatal(e.Start(":1323"))

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
