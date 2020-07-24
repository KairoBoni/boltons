package main

import (
	"fmt"
	"net/http"

	"github.com/KairoBoni/boltons/pkg/database"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db database.StoreInterface
}

func (h *Handler) getNfeAmount(c echo.Context) error {
	accessKey := c.Param("accessKey")

	amount, err := h.db.GetNfeAmount(accessKey)
	if amount == "" {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("No nfe found from the access key %s", accessKey))
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, amount)
}
