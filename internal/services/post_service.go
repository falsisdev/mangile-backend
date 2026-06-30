package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func PostToSanity(document map[string]interface{}) error {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	token := os.Getenv("SANITY_TOKEN")
	url := fmt.Sprintf("https://%s.api.sanity.io/v2021-10-21/data/mutate/production", projectID) //Burdaki %s üstteki değişkenlerden string olan ilki olduğu için projectID'yi implante etmiş olduk

	payload := map[string]interface{}{
		"mutations": []map[string]interface{}{
			{"create": document},
		},
	}

	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("[Hata]: Sanity hatası %d", resp.StatusCode)
	}

	return nil
}
