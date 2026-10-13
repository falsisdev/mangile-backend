package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
)

const (
	malAuthURL  = "https://myanimelist.net/v1/oauth2/authorize"
	malTokenURL = "https://myanimelist.net/v1/oauth2/token"
	malAPIBase  = "https://api.myanimelist.net/v2"
)

type MALOAuthTokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MALUserStatusResponse struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	MainPicture  struct {
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"main_picture"`
	NumChapters  int `json:"num_chapters"`
	MyListStatus struct {
		Status          string  `json:"status"`
		Score           float64 `json:"score"`
		NumChaptersRead float64 `json:"num_chapters_read"`
		UpdatedAt       string  `json:"updated_at"`
	} `json:"my_list_status"`
}

type MALService struct {
	httpClient *http.Client
}

func NewMALService() *MALService {
	return &MALService{
		httpClient: &http.Client{Timeout: 12 * time.Second},
	}
}

// GeneratePKCE generates code_verifier and code_challenge for MAL OAuth 2.0 PKCE
func GeneratePKCE() (verifier, challenge string, err error) {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	// Plain base64-url without padding as per PKCE RFC
	verifier = base64.RawURLEncoding.EncodeToString(bytes)
	// MAL supports plain or s256 (code_challenge_method=plain is standard for MAL, but s256 works too)
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return verifier, challenge, nil
}

func (s *MALService) GetAuthorizationURL(state, codeChallenge string) string {
	clientID := os.Getenv("MAL_CLIENT_ID")
	redirectURI := os.Getenv("MAL_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "https://mangile.vercel.app/api/sync/oauth/mal/callback"
	}

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", clientID)
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "plain")
	params.Set("redirect_uri", redirectURI)

	return fmt.Sprintf("%s?%s", malAuthURL, params.Encode())
}

func (s *MALService) ExchangeCode(ctx context.Context, code, codeVerifier string) (*MALOAuthTokenResponse, error) {
	clientID := os.Getenv("MAL_CLIENT_ID")
	clientSecret := os.Getenv("MAL_CLIENT_SECRET")
	redirectURI := os.Getenv("MAL_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "https://mangile.vercel.app/api/sync/oauth/mal/callback"
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	if clientSecret != "" {
		data.Set("client_secret", clientSecret)
	}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("code_verifier", codeVerifier)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, malTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("MAL token takasi basarisiz: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("MAL token hatasi (kod %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp MALOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (s *MALService) RefreshToken(ctx context.Context, refreshToken string) (*MALOAuthTokenResponse, error) {
	clientID := os.Getenv("MAL_CLIENT_ID")
	clientSecret := os.Getenv("MAL_CLIENT_SECRET")

	data := url.Values{}
	data.Set("client_id", clientID)
	if clientSecret != "" {
		data.Set("client_secret", clientSecret)
	}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, malTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("MAL token yenileme basarisiz: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("MAL refresh hatasi (kod %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp MALOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (s *MALService) SearchManga(ctx context.Context, query string) ([]models.TrackerSearchResult, error) {
	// First preference: Jikan API for high-detail public search without rate limit on client token
	jikanURL := fmt.Sprintf("https://api.jikan.moe/v4/manga?q=%s&limit=10", url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jikanURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var jikanRes struct {
			Data []struct {
				MalID  int    `json:"mal_id"`
				Title  string `json:"title"`
				Titles []struct {
					Type  string `json:"type"`
					Title string `json:"title"`
				} `json:"titles"`
				Images struct {
					Jpg struct {
						LargeImageURL string `json:"large_image_url"`
					} `json:"jpg"`
				} `json:"images"`
				Type     string `json:"type"`
				Chapters int    `json:"chapters"`
				Status   string `json:"status"`
				URL      string `json:"url"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&jikanRes); err == nil {
			var results []models.TrackerSearchResult
			for _, item := range jikanRes.Data {
				results = append(results, models.TrackerSearchResult{
					ID:            strconv.Itoa(item.MalID),
					Provider:      "mal",
					Title:         item.Title,
					CoverImage:    item.Images.Jpg.LargeImageURL,
					Format:        item.Type,
					Status:        item.Status,
					TotalChapters: item.Chapters,
					URL:           item.URL,
				})
			}
			return results, nil
		}
	}

	// Fallback to MAL API v2 search if client ID is configured
	clientID := os.Getenv("MAL_CLIENT_ID")
	if clientID == "" {
		return []models.TrackerSearchResult{}, nil
	}

	malSearchURL := fmt.Sprintf("%s/manga?q=%s&limit=10&fields=id,title,main_picture,status,num_chapters", malAPIBase, url.QueryEscape(query))
	req2, err := http.NewRequestWithContext(ctx, http.MethodGet, malSearchURL, nil)
	if err != nil {
		return nil, err
	}
	req2.Header.Set("X-MAL-CLIENT-ID", clientID)

	resp2, err := s.httpClient.Do(req2)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	var malRes struct {
		Data []struct {
			Node struct {
				ID          int    `json:"id"`
				Title       string `json:"title"`
				MainPicture struct {
					Medium string `json:"medium"`
					Large  string `json:"large"`
				} `json:"main_picture"`
				Status      string `json:"status"`
				NumChapters int    `json:"num_chapters"`
			} `json:"node"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&malRes); err != nil {
		return nil, err
	}

	var results []models.TrackerSearchResult
	for _, it := range malRes.Data {
		results = append(results, models.TrackerSearchResult{
			ID:            strconv.Itoa(it.Node.ID),
			Provider:      "mal",
			Title:         it.Node.Title,
			CoverImage:    it.Node.MainPicture.Large,
			Status:        it.Node.Status,
			TotalChapters: it.Node.NumChapters,
			URL:           fmt.Sprintf("https://myanimelist.net/manga/%d", it.Node.ID),
		})
	}
	return results, nil
}

func (s *MALService) GetUserMangaStatus(ctx context.Context, accessToken string, malID int) (*models.TrackerStatus, error) {
	apiURL := fmt.Sprintf("%s/manga/%d?fields=id,title,main_picture,num_chapters,my_list_status", malAPIBase, malID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("MAL get status error (%d): %s", resp.StatusCode, string(body))
	}

	var data MALUserStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	st := &models.TrackerStatus{
		Provider:        "mal",
		MatchedID:       strconv.Itoa(data.ID),
		Title:           data.Title,
		CoverImage:      data.MainPicture.Large,
		Status:          normalizeStatus(data.MyListStatus.Status),
		Score:           data.MyListStatus.Score,
		NumReadChapters: data.MyListStatus.NumChaptersRead,
		TotalChapters:   data.NumChapters,
		URL:             fmt.Sprintf("https://myanimelist.net/manga/%d", data.ID),
		IsConnected:     true,
	}
	return st, nil
}

func (s *MALService) UpdateUserMangaStatus(ctx context.Context, accessToken string, malID int, status string, chapters float64, score float64) error {
	apiURL := fmt.Sprintf("%s/manga/%d/my_list_status", malAPIBase, malID)

	form := url.Values{}
	malStatus := toMALStatus(status)
	if malStatus != "" {
		form.Set("status", malStatus)
	}
	if chapters >= 0 {
		form.Set("num_chapters_read", strconv.Itoa(int(chapters)))
	}
	if score >= 0 {
		form.Set("score", strconv.Itoa(int(score)))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("MAL guncelleme hatasi (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

func normalizeStatus(s string) string {
	switch strings.ToLower(s) {
	case "reading", "current":
		return "reading"
	case "completed":
		return "completed"
	case "on_hold", "paused":
		return "on_hold"
	case "dropped":
		return "dropped"
	case "plan_to_read", "planning":
		return "plan_to_read"
	case "repeating", "re_reading", "rereading":
		return "rereading"
	default:
		return s
	}
}

func toMALStatus(s string) string {
	switch strings.ToLower(s) {
	case "reading", "rereading":
		return "reading"
	case "completed":
		return "completed"
	case "on_hold":
		return "on_hold"
	case "dropped":
		return "dropped"
	case "plan_to_read":
		return "plan_to_read"
	default:
		return "reading"
	}
}
