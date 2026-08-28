package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	sanityAPIVersion = "v2021-10-21"
	sanityDataset    = "production"
)

// sanityClient, Sanity GROQ sorgularini calistirmak icin ortak istemcidir.
// Tum okuma sorgulari draft perspektifiyle ve varsa SANITY_TOKEN ile gonderilir;
// boylece studyoda yayinlanmamis (draft) dokumanlar da API uzerinden gorunur olur.
type sanityClient struct {
	projectID string
	token     string
	http      *http.Client
}

func newSanityClient() (*sanityClient, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("SANITY_PROJECT_ID environment degiskeni bulunamadi")
	}

	return &sanityClient{
		projectID: projectID,
		token:     os.Getenv("SANITY_TOKEN"),
		http:      &http.Client{Timeout: 15 * time.Second},
	}, nil
}

// query, verilen GROQ sorgusunu calistirir ve "result" alanini dst'e cozumler.
// params'taki degerler JSON'a kodlanip $<key> GROQ parametreleri olarak
// gonderilir; boylece kullanici girisi asla sorgu metnine karismaz.
func (c *sanityClient) query(ctx context.Context, query string, params map[string]any, dst any) error {
	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.api.sanity.io", c.projectID),
		Path:   fmt.Sprintf("/%s/data/query/%s", sanityAPIVersion, sanityDataset),
	}

	values := url.Values{}
	values.Set("query", query)
	values.Set("perspective", "drafts")
	for key, value := range params {
		encoded, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("$%s parametresi kodlanamadi: %w", key, err)
		}
		values.Set("$"+key, string(encoded))
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sanity api hatasi, durum kodu: %d", resp.StatusCode)
	}

	var wrapper struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return err
	}

	return json.Unmarshal(wrapper.Result, dst)
}
