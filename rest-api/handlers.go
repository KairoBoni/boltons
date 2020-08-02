package main

import (
	"fmt"
	"net/http"

	"github.com/KairoBoni/boltons/pkg/database"
	"github.com/labstack/echo/v4"
)

//Handler implement all methos of the REST-API
type Handler struct {
	db database.StoreInterface
}

//getNfeAmount get the total of NFe value
func (h *Handler) getNfeAmount(c echo.Context) error {
	accessKey := c.Param("accessKey")

	amount, err := h.db.GetNfeAmount(accessKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if amount == "" {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("No nfe found from the access key %s", accessKey))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf(`{"amount": "%s"}`, amount))
}
