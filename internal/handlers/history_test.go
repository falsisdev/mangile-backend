package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/falsisdev/mangile-backend/internal/handlers"
	"github.com/labstack/echo/v4"
)

func TestGetHistoryHandlerAuth(t *testing.T) {
	e := echo.New()

	t.Run("missing viewer auth returns 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/user-1/history", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-1")

		err := handlers.GetHistoryHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})

	t.Run("mismatched viewer auth returns 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/user-1/history", nil)
		req.Header.Set("X-Logto-Id", "user-2")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-1")

		err := handlers.GetHistoryHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestDeleteHistoryItemHandlerAuth(t *testing.T) {
	e := echo.New()

	t.Run("missing viewer auth returns 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/user/user-1/history/series-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id", "seriesId")
		c.SetParamValues("user-1", "series-1")

		err := handlers.DeleteHistoryItemHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestClearHistoryHandlerAuth(t *testing.T) {
	e := echo.New()

	t.Run("mismatched viewer auth returns 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/user/user-1/history", nil)
		req.Header.Set("X-Logto-Id", "another-user")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-1")

		err := handlers.ClearHistoryHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestSyncHistoryHandlerValidation(t *testing.T) {
	e := echo.New()

	t.Run("invalid json body returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/user/user-1/history/sync", strings.NewReader(`not-json`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-Logto-Id", "user-1")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("user-1")

		err := handlers.SyncHistoryHandler(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}
