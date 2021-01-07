package handlers

import (
	"encoding/json"
	"github.com/guseinovmg/imhioTest/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var articleJSON = `{"content":"Super text","tags":["MegaTag", "Super tag"]}`
var updatedArticleJSON = `{"content":"new text","tags":["New tag", "Super tag"]}`

var idStr string

func TestCreateNewArticle(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/article", strings.NewReader(articleJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, CreateNewArticle(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
	idStr = rec.Body.String()
	_, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		t.Error("id must be number")
	}

	//Check  Bad Request
	req = httptest.NewRequest(http.MethodPost, "/article", strings.NewReader(`{"content":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, CreateNewArticle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateArticle(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/article/", strings.NewReader(updatedArticleJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	if assert.NoError(t, UpdateArticle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	}

	//Check  Bad Request
	req = httptest.NewRequest(http.MethodPut, "/article", strings.NewReader(`{"content":"rrr"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	if assert.NoError(t, UpdateArticle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotEqual(t, "OK", rec.Body.String())
	}
}

func TestGetArticleById(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/article/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	if assert.NoError(t, GetArticleById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var article models.Article
		assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &article))
		assert.Equal(t, idStr, strconv.FormatUint(article.Id, 10))
	}
}

func TestGetArticlesByTag(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/article?tag=Super%20tag", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetArticlesByTag(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var articles []models.Article
		assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &articles))
	}
}

func TestDeleteArticle(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/article", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	//Assertions
	if assert.NoError(t, DeleteArticle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
