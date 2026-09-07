package handlers

import (
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

func GetUserHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Kullanıcının Logto ID'si girilmemiş."})
	}

	user, err := services.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func CreateUserHandler(c echo.Context) error {
	var body struct {
		LogtoID  string `json:"logtoId"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Avatar   string `json:"avatar"`
		Banner   string `json:"banner"`
		Bio      string `json:"bio"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: Geçersiz istek gövdesi."})
	}

	if body.LogtoID == "" || body.Username == "" || body.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "[HATA]: logtoId, username ve name alanları zorunludur."})
	}

	doc := map[string]any{
		"_type":    "auth",
		"logtoId":  body.LogtoID,
		"username": body.Username,
		"name":     body.Name,
		"roles":    []string{"004"},
		"lists":    []any{},
	}
	if body.Avatar != "" {
		doc["avatar"] = body.Avatar
	}
	if body.Banner != "" {
		doc["banner"] = body.Banner
	}
	if body.Bio != "" {
		doc["bio"] = body.Bio
	}

	if err := services.PostToSanity(c.Request().Context(), doc); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"code": 500, "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"message": "Kullanıcı başarıyla oluşturuldu.",
		"user":    doc,
	})
}

