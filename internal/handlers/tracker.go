package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/labstack/echo/v4"
)

type TrackerHandler struct {
	syncService     *services.TrackerSyncService
	malService      *services.MALService
	aniListService  *services.AniListTrackerService
	mangaDexService *services.MangaDexTrackerService
}

func NewTrackerHandler(
	syncService *services.TrackerSyncService,
	mal *services.MALService,
	aniList *services.AniListTrackerService,
	md *services.MangaDexTrackerService,
) *TrackerHandler {
	return &TrackerHandler{
		syncService:     syncService,
		malService:      mal,
		aniListService:  aniList,
		mangaDexService: md,
	}
}

func (h *TrackerHandler) GetTrackersHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	if targetLogtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "Kullanici ID'si girilmemis."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, "")
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "Bu takipci verilerine yalnizca hesap sahibi erisebilir."})
	}

	doc, err := h.syncService.GetUserTrackersDoc(c.Request().Context(), targetLogtoID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, doc)
}

func (h *TrackerHandler) GetSeriesTrackerStatusHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	seriesID := c.Param("seriesId")
	if targetLogtoID == "" || seriesID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "Eksik parametreler."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, "")
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "Yetkisiz istek."})
	}

	status, err := h.syncService.GetSeriesLiveStatus(c.Request().Context(), targetLogtoID, seriesID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, status)
}

func (h *TrackerHandler) SearchTrackerTitlesHandler(c echo.Context) error {
	provider := strings.ToLower(c.QueryParam("provider"))
	query := c.QueryParam("q")
	if provider == "" || query == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "Saglayici ve arama sorgusu zorunludur."})
	}

	ctx := c.Request().Context()
	switch provider {
	case "mal":
		results, err := h.malService.SearchManga(ctx, query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"code": 500, "message": err.Error()})
		}
		return c.JSON(http.StatusOK, results)
	case "anilist":
		results, err := h.aniListService.SearchManga(ctx, query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"code": 500, "message": err.Error()})
		}
		return c.JSON(http.StatusOK, results)
	case "mangadex":
		results, err := h.mangaDexService.SearchManga(ctx, query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"code": 500, "message": err.Error()})
		}
		return c.JSON(http.StatusOK, results)
	default:
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "Desteklenmeyen saglayici."})
	}
}

func (h *TrackerHandler) BindTrackerHandler(c echo.Context) error {
	var req models.BindTrackerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "Gecersiz istek govdesi."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, req.LogtoID)
	if viewerLogtoID == "" {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "Giris yapilmalidir."})
	}
	req.LogtoID = viewerLogtoID

	if err := h.syncService.BindTracker(c.Request().Context(), req); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "message": "Takipci eslestirmesi basariyla kaydedildi."})
}

func (h *TrackerHandler) UnbindTrackerHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	seriesID := c.Param("seriesId")
	provider := c.Param("provider")

	viewerLogtoID := getLogtoIDFromContext(c, "")
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "Yetkisiz islem."})
	}

	if err := h.syncService.UnbindTracker(c.Request().Context(), targetLogtoID, seriesID, provider); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "message": "Takipci baglantisi kaldirildi."})
}

func (h *TrackerHandler) UpdateTrackerStatusHandler(c echo.Context) error {
	var req models.UpdateTrackerStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "Gecersiz istek govdesi."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, req.LogtoID)
	if viewerLogtoID == "" {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "Giris yapilmalidir."})
	}
	req.LogtoID = viewerLogtoID

	if err := h.syncService.UpdateSingleTracker(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "message": "Durum guncellendi."})
}

func (h *TrackerHandler) SyncTrackerProgressHandler(c echo.Context) error {
	var req models.SyncProgressRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "Gecersiz istek govdesi."})
	}

	viewerLogtoID := getLogtoIDFromContext(c, req.LogtoID)
	if viewerLogtoID == "" {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "Giris yapilmalidir."})
	}
	req.LogtoID = viewerLogtoID

	// Non-blocking background sync
	h.syncService.SyncSeriesProgress(c.Request().Context(), req)

	return c.JSON(http.StatusOK, map[string]any{"success": true, "message": "Senkronizasyon arka planda baslatildi."})
}

func (h *TrackerHandler) OAuthInitiateHandler(c echo.Context) error {
	provider := strings.ToLower(c.Param("provider"))
	logtoID := c.QueryParam("logtoId")
	if logtoID == "" {
		logtoID = getLogtoIDFromContext(c, "")
	}
	if logtoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "logtoId zorunludur."})
	}

	state := fmt.Sprintf("%s:%s", provider, logtoID)

	switch provider {
	case "mal":
		verifier, challenge, err := services.GeneratePKCE()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"code": 500, "message": err.Error()})
		}
		// Include verifier in state for callback handling
		stateWithVerifier := fmt.Sprintf("%s:%s:%s", provider, logtoID, verifier)
		authURL := h.malService.GetAuthorizationURL(stateWithVerifier, challenge)
		return c.JSON(http.StatusOK, map[string]any{
			"url":          authURL,
			"codeVerifier": verifier,
		})
	case "anilist":
		authURL := h.aniListService.GetAuthorizationURL(state)
		return c.JSON(http.StatusOK, map[string]any{"url": authURL})
	case "mangadex":
		authURL := h.mangaDexService.GetAuthorizationURL(state)
		return c.JSON(http.StatusOK, map[string]any{"url": authURL})
	default:
		return c.JSON(http.StatusBadRequest, map[string]any{"code": 400, "message": "Bilinmeyen saglayici."})
	}
}

func (h *TrackerHandler) OAuthCallbackHandler(c echo.Context) error {
	provider := strings.ToLower(c.Param("provider"))
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "https://mangile.vercel.app/sync?error=missing_code")
	}

	parts := strings.Split(state, ":")
	if len(parts) < 2 {
		return c.Redirect(http.StatusTemporaryRedirect, "https://mangile.vercel.app/sync?error=invalid_state")
	}
	logtoID := parts[1]

	ctx := c.Request().Context()
	var patchSet map[string]any

	switch provider {
	case "mal":
		codeVerifier := ""
		if len(parts) >= 3 {
			codeVerifier = parts[2]
		}
		tokenResp, err := h.malService.ExchangeCode(ctx, code, codeVerifier)
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://mangile.vercel.app/sync?error=%s", err.Error()))
		}
		expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC().Format(time.RFC3339)
		patchSet = map[string]any{
			"integrations.mal": map[string]any{
				"accessToken":  tokenResp.AccessToken,
				"refreshToken": tokenResp.RefreshToken,
				"expiresAt":    expiresAt,
				"connectedAt":  time.Now().UTC().Format(time.RFC3339),
				"lastSyncAt":   time.Now().UTC().Format(time.RFC3339),
			},
		}

	case "anilist":
		tokenResp, err := h.aniListService.ExchangeCode(ctx, code)
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://mangile.vercel.app/sync?error=%s", err.Error()))
		}
		patchSet = map[string]any{
			"integrations.anilist": map[string]any{
				"accessToken": tokenResp.AccessToken,
				"connectedAt": time.Now().UTC().Format(time.RFC3339),
				"lastSyncAt":  time.Now().UTC().Format(time.RFC3339),
			},
		}

	case "mangadex":
		tokenResp, err := h.mangaDexService.ExchangeCode(ctx, code)
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://mangile.vercel.app/sync?error=%s", err.Error()))
		}
		expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC().Format(time.RFC3339)
		patchSet = map[string]any{
			"integrations.mangadex": map[string]any{
				"accessToken":  tokenResp.AccessToken,
				"refreshToken": tokenResp.RefreshToken,
				"expiresAt":    expiresAt,
				"connectedAt":  time.Now().UTC().Format(time.RFC3339),
				"lastSyncAt":   time.Now().UTC().Format(time.RFC3339),
			},
		}

	default:
		return c.Redirect(http.StatusTemporaryRedirect, "https://mangile.vercel.app/sync?error=unknown_provider")
	}

	// Update user's auth document with tokens
	if err := h.syncService.SaveUserIntegration(ctx, logtoID, patchSet); err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://mangile.vercel.app/sync?error=%s", err.Error()))
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://mangile.vercel.app/sync?connected=%s", provider))
}

func (h *TrackerHandler) DisconnectIntegrationHandler(c echo.Context) error {
	targetLogtoID := c.Param("id")
	provider := strings.ToLower(c.Param("provider"))

	viewerLogtoID := getLogtoIDFromContext(c, "")
	if viewerLogtoID == "" || viewerLogtoID != targetLogtoID {
		return c.JSON(http.StatusForbidden, map[string]any{"code": 403, "message": "Yetkisiz islem."})
	}

	if err := h.syncService.DisconnectUserIntegration(c.Request().Context(), targetLogtoID, provider); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "message": "Entegrasyon baglantisi kesildi."})
}
