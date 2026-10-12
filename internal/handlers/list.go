package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func getLogtoIDFromContext(c echo.Context, bodyLogtoID string) string {
	if bodyLogtoID != "" {
		return bodyLogtoID
	}
	headerID := c.Request().Header.Get("X-Logto-Id")
	if headerID != "" {
		return headerID
	}
	return c.QueryParam("logtoId")
}

func handleServiceError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return c.JSON(http.StatusNotFound, map[string]any{"code": 404, "message": err.Error()})
	case errors.Is(err, services.ErrUnauthorized):
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": err.Error()})
	case errors.Is(err, services.ErrListIsPrivate):
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": err.Error()})
	case errors.Is(err, services.ErrUserNotFound):
		return c.JSON(http.StatusNotFound, map[string]any{"code": 404, "message": err.Error()})
	case errors.Is(err, services.ErrInvalidRequest):
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": err.Error()})
	case errors.Is(err, services.ErrItemAlreadyExists):
		return c.JSON(http.StatusConflict, map[string]any{"code": 409, "message": err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]any{"code": 500, "message": err.Error()})
	}
}

func GetListHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Listenin ID'si girilmemiş."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, "")
	list, err := services.GetList(c.Request().Context(), id, viewerLogtoID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, list)
}

func GetPublicListsHandler(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	lists, err := services.GetPublicLists(c.Request().Context(), limit)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, lists)
}

func GetUserListsHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	if targetLogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Kullanici ID'si girilmemis."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, "")
	lists, err := services.GetUserLists(c.Request().Context(), targetLogtoID, viewerLogtoID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, lists)
}

func CreateListHandler(c echo.Context) error {
	var body models.CreateListRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Gecersiz istek govdesi."})
	}

	body.LogtoID = getLogtoIDFromContext(c, body.LogtoID)
	if body.LogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: logtoId alani zorunludur."})
	}

	if body.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: title alani zorunludur."})
	}

	list, err := services.CreateList(c.Request().Context(), body)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, list)
}

func UpdateListHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Listenin ID'si girilmemis."})
	}

	var body models.UpdateListRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Gecersiz istek govdesi."})
	}

	body.LogtoID = getLogtoIDFromContext(c, body.LogtoID)
	if body.LogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: logtoId alani zorunludur."})
	}

	list, err := services.UpdateList(c.Request().Context(), id, body)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, list)
}

func DeleteListHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Listenin ID'si girilmemis."})
	}

	var body struct {
		LogtoID string `json:"logtoId"`
	}
	_ = c.Bind(&body)

	logtoID := getLogtoIDFromContext(c, body.LogtoID)
	if logtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: logtoId alani zorunludur."})
	}

	if err := services.DeleteList(c.Request().Context(), id, logtoID); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Liste basariyla silindi.",
	})
}

func AddListItemHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Listenin ID'si girilmemis."})
	}

	var body models.AddListItemRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Gecersiz istek govdesi."})
	}

	body.LogtoID = getLogtoIDFromContext(c, body.LogtoID)
	if body.LogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: logtoId alani zorunludur."})
	}

	list, err := services.AddListItem(c.Request().Context(), id, body)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, list)
}

func RemoveListItemHandler(c echo.Context) error {
	id := c.Param("id")
	key := c.Param("key")
	if id == "" || key == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Liste ID'si ve oge key'i zorunludur."})
	}

	var body struct {
		LogtoID string `json:"logtoId"`
	}
	_ = c.Bind(&body)

	logtoID := getLogtoIDFromContext(c, body.LogtoID)
	if logtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: logtoId alani zorunludur."})
	}

	list, err := services.RemoveListItem(c.Request().Context(), id, key, logtoID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, list)
}

func UpdateListItemsHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Listenin ID'si girilmemis."})
	}

	var body models.UpdateListItemsRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Gecersiz istek govdesi."})
	}

	body.LogtoID = getLogtoIDFromContext(c, body.LogtoID)
	if body.LogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: logtoId alani zorunludur."})
	}

	list, err := services.UpdateListItems(c.Request().Context(), id, body)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, list)
}

func ToggleLikeListHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Listenin ID'si girilmemis."})
	}

	var body models.ToggleLikeListRequest
	_ = c.Bind(&body)

	logtoID := getLogtoIDFromContext(c, body.LogtoID)
	if logtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: logtoId alani zorunludur."})
	}

	list, err := services.ToggleLikeList(c.Request().Context(), id, logtoID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, list)
}

func GetLikedListsHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Kullanici ID'si girilmemis."})
	}

	lists, err := services.GetLikedLists(c.Request().Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	if lists == nil {
		lists = []models.List{}
	}

	return c.JSON(http.StatusOK, lists)
}

