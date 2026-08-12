package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func GetScan(id string) (*models.Scan, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	query := fmt.Sprintf(`*[_type == "scan" && _id == "%s"][0]{
	_id,
	_type,
	name,
	description,
	"coverImage": coverImage.asset -> url,
	"logo": logo.asset -> url,
	members,
	website
	}`, id)
	baseURL := fmt.Sprintf("https://%s.api.sanity.io/v2021-10-21/data/query/production", projectID)
	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var scanWrapper struct {
		Result models.Scan `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&scanWrapper); err != nil {
		return nil, err
	}
	return &scanWrapper.Result, nil
}
