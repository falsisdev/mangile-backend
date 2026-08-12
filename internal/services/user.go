package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func GetUser(id string) (*models.User, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	query := fmt.Sprintf(`*[_type == "auth" && logtoId == "%s"][0]{
	_id,
	_type,
	_createdAt,
	avatar,
	banner,
	bio,
	integrations,
	roles,
	logtoId,
	name,
	username,
	"lists": lists[]->{
			"id": _id,
			_type,
			title,
			createdAt,
			items,
			user->{"id": _id, logtoId, name, avatar, username},
			"likes": likes[]->{
					"id": _id,
					name,
					avatar,
					username,
					logtoId
					},
			},
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
	var userWrapper struct {
		Result models.User `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userWrapper); err != nil {
		return nil, err
	}
	return &userWrapper.Result, nil
}
