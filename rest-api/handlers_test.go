package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KairoBoni/boltons/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetNfeAmount(t *testing.T) {
	tests := []struct {
		db     *database.Mock
		status int
		body   string
	}{
		{
			db: &database.Mock{
				Amount: "120 dol",
				Err:    nil,
			},
			body:   "\"120 dol\"\n",
			status: http.StatusOK,
		},
		{
			db: &database.Mock{
				Amount: "",
				Err:    nil,
			},
			body:   "\"No nfe found from the access key blaa\"\n",
			status: http.StatusNotFound,
		},
		{
			db: &database.Mock{
				Amount: "",
				Err:    fmt.Errorf("Failed to get data from db"),
			},
			body:   "\"No nfe found from the access key blaa\"\n",
			status: http.StatusInternalServerError,
		},
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/nfe/amount/:accessKey")
	c.SetParamNames("accessKey")
	c.SetParamValues("blaa")

	for _, test := range tests {
		h := &Handler{
			db: test.db,
		}

		if assert.NoError(t, h.getNfeAmount(c)) {
			assert.Equal(t, test.status, rec.Code)
			assert.Equal(t, test.body, rec.Body.String())
		}
	}
}
