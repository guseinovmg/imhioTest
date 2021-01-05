package handlers

import (
	"context"
	"github.com/guseinovmg/imhioTest/db"
	"github.com/guseinovmg/imhioTest/models"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetArticleById(c echo.Context) error {
	id := c.Param("id")
	row := db.DB.QueryRow(context.Background(), "SELECT id,content,tags FROM articles WHERE id=$1", id)
	article := models.Article{}
	err := row.Scan(&article.Id, &article.Content, &article.Tags)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, article)
}

func GetArticleByTag(c echo.Context) error {
	tag := c.QueryParam("tag")
	res, err := db.DB.Query(context.Background(), "SELECT id,content,tags FROM articles WHERE $1=ANY(tags)", tag)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	rows := make([]models.Article, 0)
	for res.Next() {
		row := models.Article{}
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
}

func CreateNewArticle(c echo.Context) error {
	article := &models.Article{}
	if err := c.Bind(article); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err := db.DB.Exec(context.Background(), "INSERT INTO articles (content, tags) VALUES ($1, $2)", article.Content, article.Tags)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, "OK")
}

func UpdateArticle(c echo.Context) error {
	article := &models.Article{}
	id := c.Param("id")
	if err := c.Bind(article); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	commandTag, err := db.DB.Exec(context.Background(), "UPDATE articles SET content=$1, tags=$2 WHERE id=$3", article.Content, article.Tags, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if commandTag.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, "OK")
}

func DeleteArticle(c echo.Context) error {
	id := c.Param("id")
	commandTag, err := db.DB.Exec(context.Background(), "DELETE FROM articles WHERE id=$1", id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if commandTag.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.String(http.StatusOK, "OK")
}
