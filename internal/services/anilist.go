package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

const aniListAPIURL = "https://graphql.anilist.co"

// postAniListQuery, verilen GraphQL sorgusunu degiskenleriyle AniList
// API'sine gonderir ve yaniti dst icerisine cozumler.
func postAniListQuery(ctx context.Context, httpClient *http.Client, query string, variables map[string]any, dst any) error {
	requestBody, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return fmt.Errorf("graphql istek gövdesi hazırlanamadı: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, aniListAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("graphql isteği oluşturulamadı: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("anilist isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("anilist api hatası, durum kodu: %d, yanıt: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(dst)
}

// AniListService, AniList GraphQL API'si uzerinden sorgu yapar.
type AniListService struct {
	httpClient *http.Client
}

func NewAniListService() *AniListService {
	return &AniListService{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchDetails, verilen MAL kimligine karsilik gelen AniList manga
// medyasini detaylariyla dondurur.
func (s *AniListService) FetchDetails(ctx context.Context, malID int) (*models.AniListData, error) {
	var apiResponse models.AniListResponse
	if err := postAniListQuery(ctx, s.httpClient, queries.AniListMediaDetails, map[string]any{"idMal": malID}, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse.Data.Media, nil
}
