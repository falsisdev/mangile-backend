package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
)

type JikanService struct {
	httpClient *http.Client
	baseURL    string
}

func NewJikanService() *JikanService {
	return &JikanService{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.jikan.moe/v4",
	}
}

func (s *JikanService) FetchDetails(ctx context.Context, malID int) (*models.JikanData, error) {
	url := fmt.Sprintf("%s/manga/%d", s.baseURL, malID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("malID %d için Jikan kaydı bulunamadı", malID)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jikan API hatası, Durum kodu: %d", resp.StatusCode)
	}

	var apiResponse models.JikanResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse.Data, nil
}
