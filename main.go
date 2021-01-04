package main

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"os/signal"

	"github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "gopkg.in/inconshreveable/log15.v2"

	"context"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Article struct {
	Id      int64    `json:"id"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

var db *pgxpool.Pool

func main() {
	logger := log15adapter.NewLogger(log.New("module", "pgx"))

	poolConfig, err := pgxpool.ParseConfig(`postgresql://postgres:password@localhost:5432/articles`)
	if err != nil {
		log.Crit("Unable to parse DATABASE_URL", "error", err)
		os.Exit(1)
	}

	poolConfig.ConnConfig.Logger = logger

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	e := echo.New()
	e.Use(middleware.BodyLimit("2M"))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	articleGroup := e.Group("/article")

	articleGroup.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		row := db.QueryRow(context.Background(), "SELECT id,content,tags FROM articles WHERE id=$1", id)
		article := Article{}
		err = row.Scan(&article.Id, &article.Content, &article.Tags)
		if err != nil {
			if err == pgx.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound)
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}
		return c.JSON(http.StatusOK, article)
	})

	articleGroup.GET("/", func(c echo.Context) error {
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
		tag := c.QueryParam("tag")
		res, err := db.Query(context.Background(), "SELECT id,content,tags FROM articles WHERE $1=ANY(tags)", tag)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		rows := make([]Article, 0)
		for res.Next() {
			row := Article{}
			err = res.Scan(&row.Id, &row.Content, &row.Tags)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			rows = append(rows, row)
		}
		if len(rows) == 0 {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, rows)
	})

	adminGroup := e.Group("/admin")

	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == "admin" && password == "password", nil
	}))

	adminArticle := adminGroup.Group("/article")

	adminArticle.POST("", func(c echo.Context) error {
		article := &Article{}
		if err := c.Bind(article); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		_, err = db.Exec(context.Background(), "INSERT INTO articles (content, tags) VALUES ($1, $2)", article.Content, article.Tags)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, "OK")
	})

	adminArticle.PUT("/:id", func(c echo.Context) error {
		article := &Article{}
		id := c.Param("id")
		if err := c.Bind(article); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		commandTag, err := db.Exec(context.Background(), "UPDATE articles SET content=$1, tags=$2 WHERE id=$3", article.Content, article.Tags, id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if commandTag.RowsAffected() == 0 {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, "OK")
	})

	adminArticle.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")
		commandTag, err := db.Exec(context.Background(), "DELETE FROM articles WHERE id=$1", id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if commandTag.RowsAffected() == 0 {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.String(http.StatusOK, "OK")
	})

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
