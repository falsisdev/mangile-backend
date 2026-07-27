package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func GetChapter(key string, filterType string) (models.Chapter, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")

	query := fmt.Sprintf(`*[ _type == "%s" && "%s" in chapters[]._key ][0] {
		myAnimeListId,
		title,
		_type,
		"chapter": chapters[_key == "%s"][0] {
			_key,
			_type,
			chapterNumber,
			"pages": pages[] {
				"url": asset->url
			},
			source-> {
				name,
				_id
			}
		},
		"chapterKeys": chapters[]._key
	}`, filterType, key, key)

	baseURL := fmt.Sprintf("https://%s.api.sanity.io/v2021-10-21/data/query/production", projectID)

	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return models.Chapter{}, err
	}
	defer resp.Body.Close()

	var chapterWrapper struct {
		Result models.Chapter `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&chapterWrapper); err != nil {
		return models.Chapter{}, err
	}

	return chapterWrapper.Result, nil
}