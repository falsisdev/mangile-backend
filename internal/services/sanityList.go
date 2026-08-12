package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func GetSanityList(filterType string) ([]models.SanityList, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	query := fmt.Sprintf(`*[_type == 'manga' || _type == 'lightNovel'] | order(_%s desc){
		_id,
		title,
		myAnimeListId,
		_createdAt,
		_updatedAt,
		_type,
		tags,
		"bannerImage": bannerImage.asset->url,
		"coverImage": coverImage.asset->url
	}`, filterType)
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
	var sanityListWrapper struct {
		Result []models.SanityList `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sanityListWrapper); err != nil {
		return nil, err
	}
	return sanityListWrapper.Result, nil
}
