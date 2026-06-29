package handlers

import (
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func GetLightNovelHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Serinin ID'si girilmemiş."})
	}

	lightNovel, err := services.GetLightNovel(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, lightNovel)
}
