package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type SanityMutationResultItem struct {
	ID        string `json:"id"`
	Operation string `json:"operation"`
}

type SanityMutationResponse struct {
	TransactionID string                     `json:"transactionId"`
	Results       []SanityMutationResultItem `json:"results"`
}

func MutateSanity(ctx context.Context, mutations []map[string]any) (*SanityMutationResponse, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("SANITY_PROJECT_ID environment degiskeni bulunamadi")
	}

	token := os.Getenv("SANITY_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("SANITY_TOKEN environment degiskeni bulunamadi")
	}

	payload := map[string]any{
		"mutations": mutations,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("mutasyon govdesi hazirlanamadi: %w", err)
	}

	endpoint := fmt.Sprintf("https://%s.api.sanity.io/%s/data/mutate/%s", projectID, sanityAPIVersion, sanityDataset)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("yanit govdesi okunamadi: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sanity mutasyon hatasi (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var mutationResp SanityMutationResponse
	if err := json.Unmarshal(bodyBytes, &mutationResp); err != nil {
		return nil, fmt.Errorf("mutasyon yaniti cozumlenemedi: %w", err)
	}

	return &mutationResp, nil
}
