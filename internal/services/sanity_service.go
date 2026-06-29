package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func GetArticle(slug string) (*models.Article, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")

	query := fmt.Sprintf(`*[_type == 'articles' && slug.current == "%s"][0]`, slug)

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

func GetUser(id string) (*models.User, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")

	query := fmt.Sprintf(`*[_type == "auth" && logtoId == "%s"][0]{
	_id,
	_type,
	title,
	_createdAt,
	avatar,
	"banner": banner.asset -> url,
	bio,
	favoriteChapters,
	favoriteTitles,
	favoriteTitle,
	gender,
	logtoId,
	name,
	username,
    "lists": lists[]->{
             _id,
			_type,
			title,
			createdAt,
			items,
			user->{_id, logtoId, name, avatar, username},
			"likes": likes[]->{
					_id,
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

func GetList(id string) (*models.List, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")

	query := fmt.Sprintf(`*[_type == "lists" && _id == "%s"][0]{
	_id,
	_type,
	title,
	createdAt,
	items,
	user->{_id, logtoId, name, avatar, username},
    "likes": likes[]->{
              _id,
              name,
              avatar,
              username,
              logtoId
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

	var listWrapper struct {
		Result models.List `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&listWrapper); err != nil {
		return nil, err
	}

	return &listWrapper.Result, nil
}

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
func GetLightNovel(id string) (*models.LightNovel, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")

	query := fmt.Sprintf(`*[_type == "lightNovel" && myAnimeListId == %s][0]{
    _id, 
    title, 
    description,
    myAnimeListId,
	_createdAt,
	_updatedAt,
	_type,
	tags,
    "bannerImage": bannerImage.asset->url,
	"coverImage": coverImage.asset->url,
	"chapters": chapters[]{
		chapterNumber,
		title,
		"content": content
	},
	notes
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

	var lightNovelWrapper struct {
		Result models.LightNovel `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&lightNovelWrapper); err != nil {
		return nil, err
	}

	return &lightNovelWrapper.Result, nil
}

func GetManga(id string) (*models.Manga, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")

	query := fmt.Sprintf(`*[_type == "manga" && myAnimeListId == %s][0]{
    _id, 
    title, 
    description,
    myAnimeListId,
	_createdAt,
	_updatedAt,
	_type,
	tags,
    "bannerImage": bannerImage.asset->url,
	"coverImage": coverImage.asset->url,
	"chapters": chapters[]{
		chapterNumber,
		title,
		"pages": pages[]{
				"asset": {
						"url": asset->url
				}
			}
		},
	notes
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

	var mangaWrapper struct {
		Result models.Manga `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mangaWrapper); err != nil {
		return nil, err
	}

	return &mangaWrapper.Result, nil
}

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
