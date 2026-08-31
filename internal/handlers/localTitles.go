package handlers

import (
	"net/http"
	"strconv"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

// GetLocalTitlesHandler, /api/localTitles endpoint'ini işler.
func GetLocalTitlesHandler(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	titleType := c.QueryParam("type")
	tag := c.QueryParam("tag")
	status := c.QueryParam("status")
	sort := c.QueryParam("sort")

	filter := services.LocalTitlesFilter{
		Search: search,
		Type:   titleType,
		Tag:    tag,
		Status: status,
		Sort:   sort,
		Page:   page,
		Limit:  limit,
	}

	response, err := services.GetLocalTitles(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}
