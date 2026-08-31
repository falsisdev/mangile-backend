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
	page := c.QueryParam("page")
	searchQuery := c.QueryParam("query")
	if searchQuery == "" {
		searchQuery = c.QueryParam("search")
	}
	sortParam := c.QueryParam("sort")
	formatParam := c.QueryParam("format")
	genreParam := c.QueryParam("genre")
	statusParam := c.QueryParam("status")

	parsedLimit := 50
	var err error
	if limit != "" {
		parsedLimit, err = strconv.Atoi(limit)
		if err != nil || parsedLimit < 1 {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"code":    400,
				"message": "[HATA]: limit parametresi geçerli bir pozitif sayı olmalıdır.",
			})
		}
	}

	parsedPage := 1
	if page != "" {
		parsedPage, err = strconv.Atoi(page)
		if err != nil || parsedPage < 1 {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"code":    400,
				"message": "[HATA]: page parametresi geçerli bir pozitif sayı olmalıdır.",
			})
		}
	}

	if filterType == "" && sortParam == "" {
		filterType = "POPULAR"
	}

	mangaList, err := services.GetMangaList(c.Request().Context(), filterType, parsedLimit, parsedPage, searchQuery, sortParam, formatParam, genreParam, statusParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "[HATA]: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"code":  200,
		"data":  mangaList,
		"page":  parsedPage,
		"limit": parsedLimit,
	})
}
