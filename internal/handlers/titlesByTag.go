package handlers

import (
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func GetTitlesByTagHandler(c echo.Context) error {
	tag := c.QueryParam("tag")
	if tag == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: tag parametresi girilmemiş."})
	}

	titles, err := services.GetTitlesByTag(c.Request().Context(), tag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, titles)
}
