package handlers

import (
	"errors"
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func GetChapterHandler(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: id parametresi girilmemiş."})
	}

	chapter, err := services.GetChapter(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrChapterNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{"code": 404, "message": "[HATA]: Bölüm bulunamadı."})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, chapter)
}
