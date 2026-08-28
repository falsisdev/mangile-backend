package handlers

import (
	"net/http"
	"strconv"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

type TitleHandler struct {
	worker *services.TitleDetailsWorker
}

func NewTitleHandler(worker *services.TitleDetailsWorker) *TitleHandler {
	return &TitleHandler{
		worker: worker,
	}
}

func (h *TitleHandler) GetTitleDetails(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.QueryParam("id")
	malIDStr := c.QueryParam("mal_id")

	if id == "" && malIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "En az bir arama parametresi ('id' veya 'mal_id') sağlanmalıdır",
		})
	}

	var malID int
	var err error

	if malIDStr != "" {
		malID, err = strconv.Atoi(malIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Geçersiz 'mal_id' formatı, bir sayı olmalıdır",
			})
		}
	}

	completeTitle, err := h.worker.GetCompleteTitle(ctx, id, malID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Detaylar getirilirken bir hata oluştu: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, completeTitle)
}
