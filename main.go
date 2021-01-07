package main

import (
	"context"
	"github.com/guseinovmg/imhioTest/handlers"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	articleGroup := e.Group("/article")

	articleGroup.GET("/:id", handlers.GetArticleById)

	articleGroup.GET("", handlers.GetArticlesByTag, handlers.SetTokenAndCounter)

	articleGroup.POST("", handlers.CreateNewArticle)

	articleGroup.PUT("/:id", handlers.UpdateArticle)

	articleGroup.DELETE("/:id", handlers.DeleteArticle)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

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
