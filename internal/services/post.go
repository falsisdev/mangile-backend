package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// PostToSanity, verilen dokumani Sanity'ye yeni dokuman olarak olusturur.
func PostToSanity(ctx context.Context, document map[string]any) error {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("SANITY_PROJECT_ID environment degiskeni bulunamadi")
	}

	token := os.Getenv("SANITY_TOKEN")
	if token == "" {
		return fmt.Errorf("SANITY_TOKEN environment degiskeni bulunamadi")
	}

	payload := map[string]any{
		"mutations": []map[string]any{
			{"create": document},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mutasyon gövdesi hazırlanamadı: %w", err)
	}

	endpoint := fmt.Sprintf("https://%s.api.sanity.io/%s/data/mutate/%s", projectID, sanityAPIVersion, sanityDataset)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[Hata]: Sanity hatası %d", resp.StatusCode)
	}

	return nil
}
