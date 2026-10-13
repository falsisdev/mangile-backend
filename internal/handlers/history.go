package handlers

import (
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func GetHistoryHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	if targetLogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Kullanici ID'si girilmemis."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, "")
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "[HATA]: Okuma gecmisine yalnizca hesap sahibi erisebilir."})
	}

	history, err := services.GetReadingHistory(c.Request().Context(), targetLogtoID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, history)
}

func UpsertHistoryHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	if targetLogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Kullanici ID'si girilmemis."})
	}

	var req models.UpsertHistoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Gecersiz istek govdesi."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, req.LogtoID)
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "[HATA]: Bu islem icin yetkiniz bulunmuyor."})
	}

	history, err := services.UpsertReadingHistoryItem(c.Request().Context(), targetLogtoID, req.Item)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, history)
}

func DeleteHistoryItemHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	seriesID := c.Param("seriesId")
	if targetLogtoID == "" || seriesID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Eksik parametreler."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, "")
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "[HATA]: Bu islem icin yetkiniz bulunmuyor."})
	}

	if err := services.RemoveReadingHistoryItem(c.Request().Context(), targetLogtoID, seriesID); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "message": "Seri okuma gecmisinden silindi."})
}

func ClearHistoryHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	if targetLogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Kullanici ID'si girilmemis."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, "")
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "[HATA]: Bu islem icin yetkiniz bulunmuyor."})
	}

	if err := services.ClearReadingHistory(c.Request().Context(), targetLogtoID); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "message": "Okuma gecmisi tamamen temizlendi."})
}

func SyncHistoryHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	if targetLogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Kullanici ID'si girilmemis."})
	}

	var req models.BatchSyncHistoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Gecersiz istek govdesi."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, req.LogtoID)
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "[HATA]: Bu islem icin yetkiniz bulunmuyor."})
	}

	history, err := services.SyncReadingHistory(c.Request().Context(), targetLogtoID, req.Items)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, history)
}
