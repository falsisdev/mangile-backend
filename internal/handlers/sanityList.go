package handlers

import (
	"errors"
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func GetSanityListHandler(c echo.Context) error {
	filterType := c.QueryParam("filterType")
	if filterType == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: filterType parametresi girilmemiş."})
	}
	sanityList, err := services.GetSanityList(c.Request().Context(), filterType)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSortField) {
			return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: filterType yalnizca 'createdAt' veya 'updatedAt' olabilir."})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, sanityList)
}
