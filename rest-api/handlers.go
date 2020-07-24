package main

import (
	"net/http"

	"github.com/KairoBoni/boltons/pkg/database"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *database.Store
}

func (h *Handler) getNfeTotal(c echo.Context) error {
	accessKey := c.Param("accessKey")

	total, err := h.db.GetNfeTotal(accessKey)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, total)
}
