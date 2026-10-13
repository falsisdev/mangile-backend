package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
)

const (
	mangaDexAPIBase   = "https://api.mangadex.org"
	mangaDexAuthBase  = "https://auth.mangadex.org/realms/mangadex/protocol/openid-connect"
	mangaDexUploadCDN = "https://uploads.mangadex.org/covers"
)

var uuidRegex = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

type MangaDexOAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type MangaDexTrackerService struct {
	httpClient *http.Client
}

func NewMangaDexTrackerService() *MangaDexTrackerService {
	return &MangaDexTrackerService{
		httpClient: &http.Client{Timeout: 12 * time.Second},
	}
}

func (s *MangaDexTrackerService) GetAuthorizationURL(state string) string {
	clientID := os.Getenv("MANGADEX_CLIENT_ID")
	redirectURI := os.Getenv("MANGADEX_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "https://mangile.vercel.app/api/sync/oauth/mangadex/callback"
	}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "openid")
	params.Set("state", state)

	return fmt.Sprintf("%s/auth?%s", mangaDexAuthBase, params.Encode())
}

func (s *MangaDexTrackerService) ExchangeCode(ctx context.Context, code string) (*MangaDexOAuthTokenResponse, error) {
	clientID := os.Getenv("MANGADEX_CLIENT_ID")
	clientSecret := os.Getenv("MANGADEX_CLIENT_SECRET")
	redirectURI := os.Getenv("MANGADEX_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "https://mangile.vercel.app/api/sync/oauth/mangadex/callback"
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mangaDexAuthBase+"/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("MangaDex token takasi basarisiz: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("MangaDex token hatasi (kod %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp MangaDexOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (s *MangaDexTrackerService) RefreshToken(ctx context.Context, refreshToken string) (*MangaDexOAuthTokenResponse, error) {
	clientID := os.Getenv("MANGADEX_CLIENT_ID")
	clientSecret := os.Getenv("MANGADEX_CLIENT_SECRET")

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mangaDexAuthBase+"/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("MangaDex token yenileme basarisiz: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("MangaDex refresh hatasi (kod %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp MangaDexOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (s *MangaDexTrackerService) SearchManga(ctx context.Context, queryOrURL string) ([]models.TrackerSearchResult, error) {
	trimmed := strings.TrimSpace(queryOrURL)

	// Check if input is a direct UUID or URL containing UUID
	if uuidMatch := uuidRegex.FindString(trimmed); uuidMatch != "" {
		single, err := s.getMangaByUUID(ctx, uuidMatch)
		if err == nil && single != nil {
			return []models.TrackerSearchResult{*single}, nil
		}
	}

	apiURL := fmt.Sprintf("%s/manga?title=%s&limit=10&includes[]=cover_art", mangaDexAPIBase, url.QueryEscape(trimmed))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []models.TrackerSearchResult{}, nil
	}

	var mdRes struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Title       map[string]string   `json:"title"`
				AltTitles   []map[string]string `json:"altTitles"`
				Status      string              `json:"status"`
				Year        int                 `json:"year"`
				LastChapter string              `json:"lastChapter"`
			} `json:"attributes"`
			Relationships []struct {
				ID         string `json:"id"`
				Type       string `json:"type"`
				Attributes struct {
					FileName string `json:"fileName"`
				} `json:"attributes"`
			} `json:"relationships"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mdRes); err != nil {
		return nil, err
	}

	var results []models.TrackerSearchResult
	for _, item := range mdRes.Data {
		title := extractMDTitle(item.Attributes.Title, item.Attributes.AltTitles)
		coverFileName := ""
		for _, rel := range item.Relationships {
			if rel.Type == "cover_art" {
				coverFileName = rel.Attributes.FileName
				break
			}
		}

		coverURL := ""
		if coverFileName != "" {
			coverURL = fmt.Sprintf("%s/%s/%s.256.jpg", mangaDexUploadCDN, item.ID, coverFileName)
		}

		totalCh := 0
		if item.Attributes.LastChapter != "" {
			var parsed int
			fmt.Sscanf(item.Attributes.LastChapter, "%d", &parsed)
			totalCh = parsed
		}

		results = append(results, models.TrackerSearchResult{
			ID:            item.ID,
			Provider:      "mangadex",
			Title:         title,
			CoverImage:    coverURL,
			Status:        item.Attributes.Status,
			TotalChapters: totalCh,
			Year:          item.Attributes.Year,
			URL:           fmt.Sprintf("https://mangadex.org/title/%s", item.ID),
		})
	}

	return results, nil
}

func (s *MangaDexTrackerService) getMangaByUUID(ctx context.Context, mangaUUID string) (*models.TrackerSearchResult, error) {
	apiURL := fmt.Sprintf("%s/manga/%s?includes[]=cover_art", mangaDexAPIBase, mangaUUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MangaDex manga bulunamadi: %s", mangaUUID)
	}

	var mdRes struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				Title     map[string]string   `json:"title"`
				AltTitles []map[string]string `json:"altTitles"`
				Status    string              `json:"status"`
				Year      int                 `json:"year"`
			} `json:"attributes"`
			Relationships []struct {
				ID         string `json:"id"`
				Type       string `json:"type"`
				Attributes struct {
					FileName string `json:"fileName"`
				} `json:"attributes"`
			} `json:"relationships"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mdRes); err != nil {
		return nil, err
	}

	item := mdRes.Data
	title := extractMDTitle(item.Attributes.Title, item.Attributes.AltTitles)
	coverFileName := ""
	for _, rel := range item.Relationships {
		if rel.Type == "cover_art" {
			coverFileName = rel.Attributes.FileName
			break
		}
	}

	coverURL := ""
	if coverFileName != "" {
		coverURL = fmt.Sprintf("%s/%s/%s.256.jpg", mangaDexUploadCDN, item.ID, coverFileName)
	}

	return &models.TrackerSearchResult{
		ID:         item.ID,
		Provider:   "mangadex",
		Title:      title,
		CoverImage: coverURL,
		Status:     item.Attributes.Status,
		Year:       item.Attributes.Year,
		URL:        fmt.Sprintf("https://mangadex.org/title/%s", item.ID),
	}, nil
}

func (s *MangaDexTrackerService) GetUserMangaStatus(ctx context.Context, accessToken string, mangaUUID string) (*models.TrackerStatus, error) {
	// 1. Fetch Manga Metadata (Public)
	meta, err := s.getMangaByUUID(ctx, mangaUUID)
	if err != nil {
		return nil, err
	}

	status := "plan_to_read"
	var score float64

	if accessToken != "" {
		// 2. Fetch Reading Status
		statusURL := fmt.Sprintf("%s/manga/%s/status", mangaDexAPIBase, mangaUUID)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, statusURL, nil)
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+accessToken)
			if resp, err := s.httpClient.Do(req); err == nil && resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()
				var stRes struct {
					Status string `json:"status"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&stRes); err == nil && stRes.Status != "" {
					status = normalizeStatus(stRes.Status)
				}
			}
		}

		// 3. Fetch Rating
		ratingURL := fmt.Sprintf("%s/rating?manga[]=%s", mangaDexAPIBase, mangaUUID)
		req2, err := http.NewRequestWithContext(ctx, http.MethodGet, ratingURL, nil)
		if err == nil {
			req2.Header.Set("Authorization", "Bearer "+accessToken)
			if resp2, err := s.httpClient.Do(req2); err == nil && resp2.StatusCode == http.StatusOK {
				defer resp2.Body.Close()
				var ratRes struct {
					Ratings map[string]struct {
						Rating float64 `json:"rating"`
					} `json:"ratings"`
				}
				if err := json.NewDecoder(resp2.Body).Decode(&ratRes); err == nil {
					if r, ok := ratRes.Ratings[mangaUUID]; ok {
						score = r.Rating
					}
				}
			}
		}
	}

	return &models.TrackerStatus{
		Provider:        "mangadex",
		MatchedID:       mangaUUID,
		Title:           meta.Title,
		CoverImage:      meta.CoverImage,
		Status:          status,
		Score:           score,
		NumReadChapters: 0,
		TotalChapters:   meta.TotalChapters,
		URL:             meta.URL,
		IsConnected:     accessToken != "",
	}, nil
}

func (s *MangaDexTrackerService) UpdateUserMangaStatus(ctx context.Context, accessToken string, mangaUUID string, status string, score float64) error {
	if accessToken == "" {
		return fmt.Errorf("MangaDex erisim tokeni bulunamadi")
	}

	// 1. Update Reading Status
	mdStatus := toMangaDexStatus(status)
	if mdStatus != "" {
		statusURL := fmt.Sprintf("%s/manga/%s/status", mangaDexAPIBase, mangaUUID)
		body, _ := json.Marshal(map[string]string{"status": mdStatus})
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, statusURL, bytes.NewBuffer(body))
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Header.Set("Content-Type", "application/json")
			resp, err := s.httpClient.Do(req)
			if err == nil {
				resp.Body.Close()
			}
		}
	}

	// 2. Update Rating if provided
	if score > 0 {
		ratingURL := fmt.Sprintf("%s/rating/%s", mangaDexAPIBase, mangaUUID)
		body, _ := json.Marshal(map[string]int{"rating": int(score)})
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, ratingURL, bytes.NewBuffer(body))
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Header.Set("Content-Type", "application/json")
			resp, err := s.httpClient.Do(req)
			if err == nil {
				resp.Body.Close()
			}
		}
	}

	return nil
}

func toMangaDexStatus(s string) string {
	switch strings.ToLower(s) {
	case "reading":
		return "reading"
	case "completed":
		return "completed"
	case "on_hold":
		return "on_hold"
	case "dropped":
		return "dropped"
	case "plan_to_read":
		return "plan_to_read"
	case "rereading":
		return "re_reading"
	default:
		return "reading"
	}
}

func extractMDTitle(titles map[string]string, altTitles []map[string]string) string {
	if t, ok := titles["en"]; ok && t != "" {
		return t
	}
	if t, ok := titles["ja-ro"]; ok && t != "" {
		return t
	}
	for _, val := range titles {
		if val != "" {
			return val
		}
	}
	for _, alt := range altTitles {
		if t, ok := alt["en"]; ok && t != "" {
			return t
		}
	}
	return "İsimsiz Manga"
}
