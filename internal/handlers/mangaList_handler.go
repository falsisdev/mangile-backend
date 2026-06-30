package handlers

import (
	"net/http"
	"strconv"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func GetMangaListHandler(c echo.Context) error {
	filterType := c.QueryParam("filterType")
	limit := c.QueryParam("limit")
	if filterType == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: filterType parametresi girilmemiş."})
	} else if limit == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: limit parametresi girilmemiş."})
	}

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil || parsedLimit < 1 {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: limit parametresi geçerli bir pozitif sayı olmalıdır."})
	}

	mangaList, err := services.GetMangaList(filterType, parsedLimit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mangaList)
}
