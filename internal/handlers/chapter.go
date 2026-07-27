package handlers

import (
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func GetChapterHandler(c echo.Context) error {
	filterType := c.QueryParam("filterType")
	key := c.QueryParam("key")
	if filterType == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: filterType parametresi girilmemiş."})
	}
	if key == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: key parametresi girilmemiş."})
	}
	chapter, err := services.GetChapter(key, filterType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, chapter)
}
