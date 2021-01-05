package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var articleJSON = `{"content":"Super text","tags":["MegaTag", "Super tag"]}`

func TestCreateNewArticle(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/article", strings.NewReader(articleJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//Assertions
	if assert.NoError(t, CreateNewArticle(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	}

	//Check  Bad Request
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"content":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, CreateNewArticle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
