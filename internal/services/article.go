package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func GetArticle(slug string) (*models.Article, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	query := fmt.Sprintf(`*[_type == 'articles' && slug.current == "%s"][0]{..., "id": _id}`, slug)
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
	var articleWrapper struct {
		Result models.Article `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&articleWrapper); err != nil {
		return nil, err
	}
	return &articleWrapper.Result, nil
}