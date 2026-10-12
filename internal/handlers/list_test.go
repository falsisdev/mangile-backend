package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/falsisdev/mangile-backend/internal/handlers"
	"github.com/labstack/echo/v4"
)

func TestCreateListHandlerValidation(t *testing.T) {
	e := echo.New()

	t.Run("missing logtoId returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/list", strings.NewReader(`{"title":"Listem"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateListHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("missing title returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/list", strings.NewReader(`{"logtoId":"user-123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateListHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestDeleteListHandlerValidation(t *testing.T) {
	e := echo.New()

	t.Run("missing logtoId returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/list/123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("123")

		err := handlers.DeleteListHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestAddListItemHandlerValidation(t *testing.T) {
	e := echo.New()

	t.Run("missing logtoId returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/list/123/items", strings.NewReader(`{"referenceId":"manga-1"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("123")

		err := handlers.AddListItemHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}
