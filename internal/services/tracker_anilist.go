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
	"strconv"
	"strings"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
)

const (
	aniListAuthURL    = "https://anilist.co/api/v2/oauth/authorize"
	aniListTokenURL   = "https://anilist.co/api/v2/oauth/token"
	aniListGraphQLURL = "https://graphql.anilist.co"
)

type AniListOAuthTokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

type AniListTrackerService struct {
	httpClient *http.Client
}

func NewAniListTrackerService() *AniListTrackerService {
	return &AniListTrackerService{
		httpClient: &http.Client{Timeout: 12 * time.Second},
	}
}

func (s *AniListTrackerService) GetAuthorizationURL(state string) string {
	clientID := os.Getenv("ANILIST_CLIENT_ID")
	redirectURI := os.Getenv("ANILIST_REDIRECT_URI")
	if redirectURI == "" {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		redirectURI = fmt.Sprintf("%s/api/integrations/anilist-callback", frontendURL)
	}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("state", state)

	return fmt.Sprintf("%s?%s", aniListAuthURL, params.Encode())
}

func (s *AniListTrackerService) ExchangeCode(ctx context.Context, code string) (*AniListOAuthTokenResponse, error) {
	clientID := os.Getenv("ANILIST_CLIENT_ID")
	clientSecret := os.Getenv("ANILIST_CLIENT_SECRET")
	redirectURI := os.Getenv("ANILIST_REDIRECT_URI")
	if redirectURI == "" {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		redirectURI = fmt.Sprintf("%s/api/integrations/anilist-callback", frontendURL)
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, aniListTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("AniList token takasi basarisiz: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("AniList token hatasi (kod %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp AniListOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (s *AniListTrackerService) SearchManga(ctx context.Context, query string) ([]models.TrackerSearchResult, error) {
	const searchGraphQL = `
query ($search: String) {
  Page(page: 1, perPage: 10) {
    media(search: $search, type: MANGA) {
      id
      title {
        romaji
        english
        native
      }
      coverImage {
        extraLarge
      }
      format
      status
      chapters
      startDate {
        year
      }
      siteUrl
    }
  }
}`
	var result struct {
		Data struct {
			Page struct {
				Media []struct {
					ID    int `json:"id"`
					Title struct {
						Romaji  string `json:"romaji"`
						English string `json:"english"`
					} `json:"title"`
					CoverImage struct {
						ExtraLarge string `json:"extraLarge"`
					} `json:"coverImage"`
					Format    string `json:"format"`
					Status    string `json:"status"`
					Chapters  int    `json:"chapters"`
					StartDate struct {
						Year int `json:"year"`
					} `json:"startDate"`
					SiteURL string `json:"siteUrl"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := postAniListQuery(ctx, s.httpClient, searchGraphQL, map[string]any{"search": query}, &result); err != nil {
		return nil, err
	}

	var out []models.TrackerSearchResult
	for _, m := range result.Data.Page.Media {
		title := m.Title.Romaji
		if title == "" {
			title = m.Title.English
		}
		out = append(out, models.TrackerSearchResult{
			ID:            strconv.Itoa(m.ID),
			Provider:      "anilist",
			Title:         title,
			RomajiTitle:   m.Title.Romaji,
			CoverImage:    m.CoverImage.ExtraLarge,
			Format:        m.Format,
			Status:        m.Status,
			TotalChapters: m.Chapters,
			Year:          m.StartDate.Year,
			URL:           m.SiteURL,
		})
	}
	return out, nil
}

func (s *AniListTrackerService) GetUserMangaStatus(ctx context.Context, accessToken string, aniListID int) (*models.TrackerStatus, error) {
	const statusGraphQL = `
query ($mediaId: Int) {
  Media(id: $mediaId, type: MANGA) {
    id
    title {
      romaji
      english
    }
    coverImage {
      extraLarge
    }
    chapters
    siteUrl
    mediaListEntry {
      id
      status
      score(format: POINT_10_DECIMAL)
      progress
    }
  }
}`

	reqBody, err := json.Marshal(map[string]any{
		"query": statusGraphQL,
		"variables": map[string]any{
			"mediaId": aniListID,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, aniListGraphQLURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiRes struct {
		Data struct {
			Media struct {
				ID    int `json:"id"`
				Title struct {
					Romaji  string `json:"romaji"`
					English string `json:"english"`
				} `json:"title"`
				CoverImage struct {
					ExtraLarge string `json:"extraLarge"`
				} `json:"coverImage"`
				Chapters       int    `json:"chapters"`
				SiteURL        string `json:"siteUrl"`
				MediaListEntry *struct {
					ID       int     `json:"id"`
					Status   string  `json:"status"`
					Score    float64 `json:"score"`
					Progress float64 `json:"progress"`
				} `json:"mediaListEntry"`
			} `json:"Media"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return nil, err
	}

	m := apiRes.Data.Media
	title := m.Title.Romaji
	if title == "" {
		title = m.Title.English
	}

	status := "plan_to_read"
	var score float64
	var progress float64

	if m.MediaListEntry != nil {
		status = normalizeStatus(m.MediaListEntry.Status)
		score = m.MediaListEntry.Score
		progress = m.MediaListEntry.Progress
	}

	return &models.TrackerStatus{
		Provider:        "anilist",
		MatchedID:       strconv.Itoa(m.ID),
		Title:           title,
		CoverImage:      m.CoverImage.ExtraLarge,
		Status:          status,
		Score:           score,
		NumReadChapters: progress,
		TotalChapters:   m.Chapters,
		URL:             m.SiteURL,
		IsConnected:     accessToken != "",
	}, nil
}

func (s *AniListTrackerService) UpdateUserMangaStatus(ctx context.Context, accessToken string, aniListID int, status string, chapters float64, score float64) error {
	const mutateGraphQL = `
mutation ($mediaId: Int, $status: MediaListStatus, $progress: Int, $score: Float) {
  SaveMediaListEntry(mediaId: $mediaId, status: $status, progress: $progress, score: $score) {
    id
    status
    progress
    score
  }
}`

	variables := map[string]any{
		"mediaId": aniListID,
	}
	if status != "" {
		variables["status"] = toAniListStatus(status)
	}
	if chapters >= 0 {
		variables["progress"] = int(chapters)
	}
	if score >= 0 {
		variables["score"] = score
	}

	reqBody, err := json.Marshal(map[string]any{
		"query":     mutateGraphQL,
		"variables": variables,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, aniListGraphQLURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("AniList guncelleme hatasi (kod %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

func toAniListStatus(s string) string {
	switch strings.ToLower(s) {
	case "reading":
		return "CURRENT"
	case "completed":
		return "COMPLETED"
	case "on_hold":
		return "PAUSED"
	case "dropped":
		return "DROPPED"
	case "plan_to_read":
		return "PLANNING"
	case "rereading":
		return "REPEATING"
	default:
		return "CURRENT"
	}
}
