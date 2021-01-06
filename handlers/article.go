package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/guseinovmg/imhioTest/db"
	"github.com/guseinovmg/imhioTest/models"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func SetTokenAndCounter(next echo.HandlerFunc) echo.HandlerFunc {
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

func GetArticleById(c echo.Context) error {
	id := c.Param("id")
	row := db.DB.QueryRow(context.Background(), "SELECT id,content,tags FROM articles WHERE id=$1", id)
	article := models.Article{}
	err := row.Scan(&article.Id, &article.Content, &article.Tags)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.String(http.StatusNotFound, echo.ErrNotFound.Error())
		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, article)
}

func GetArticlesByTag(c echo.Context) error {
	tag := c.QueryParam("tag")
	if tag == "" {
		return c.String(http.StatusBadRequest, "Parameter tag is empty")
	}
	var offset uint64 = 0
	var limit uint64 = 10
	offsetStr := c.QueryParam("offset")
	var err error

	if offsetStr != "" {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Parameter offset is not number")
		}
	}

	limitStr := c.QueryParam("limit")
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Parameter offset is not number")
		}
	}

	res, err := db.DB.Query(context.Background(), "SELECT id,content,tags FROM articles WHERE $1=ANY(tags) LIMIT $2 OFFSET $3", tag, limit, offset)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	rows := make([]models.Article, 0)
	for res.Next() {
		row := models.Article{}
		err = res.Scan(&row.Id, &row.Content, &row.Tags)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return c.JSON(http.StatusNotFound, rows)
	}
	return c.JSON(http.StatusOK, rows)
}

func CreateNewArticle(c echo.Context) error {
	article := &models.Article{}
	if err := c.Bind(article); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if article.Content == "" {
		return c.String(http.StatusBadRequest, "Content is empty")
	}
	row := db.DB.QueryRow(context.Background(), "INSERT INTO articles (content, tags) VALUES ($1, $2) RETURNING id", article.Content, article.Tags)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusCreated, strconv.Itoa(id))
}

func UpdateArticle(c echo.Context) error {
	article := &models.Article{}
	id := c.Param("id")
	if err := c.Bind(article); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if article.Content == "" {
		return c.String(http.StatusBadRequest, "Content is empty")
	}
	if article.Tags == nil {
		return c.String(http.StatusBadRequest, "Tags is empty")
	}
	commandTag, err := db.DB.Exec(context.Background(), "UPDATE articles SET content=$1, tags=$2 WHERE id=$3", article.Content, article.Tags, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if commandTag.RowsAffected() == 0 {
		return c.String(http.StatusNotFound, echo.ErrNotFound.Error())
	}
	return c.String(http.StatusOK, "OK")
}

func DeleteArticle(c echo.Context) error {
	id := c.Param("id")
	commandTag, err := db.DB.Exec(context.Background(), "DELETE FROM articles WHERE id=$1", id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if commandTag.RowsAffected() == 0 {
		return c.String(http.StatusNotFound, echo.ErrNotFound.Error())
	}
	return c.String(http.StatusOK, "OK")
}
