package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func fetchJikanMangaData(malID int) (*models.JikanMangaResponse, error) {
	url := fmt.Sprintf("https://api.jikan.moe/v4/manga/%d", malID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var jikanResp models.JikanMangaResponse
	if err := json.NewDecoder(resp.Body).Decode(&jikanResp); err != nil {
		return nil, err
	}
	return &jikanResp, nil
}

func fetchAniListMangaData(malID int) (*models.AniListMangaResponse, error) {
	jsonData := map[string]interface{}{
		"query": fmt.Sprintf(`{
			Media(idMal: %d, type: MANGA) {
				id
				title {
					romaji
					english
					native
				}
				trending
				averageScore
				bannerImage
				coverImage {
					large
				}
				description
				tags {
					name
				}
				relations {
					edges {
						id
						relationType
						node {
							coverImage {
								extraLarge
							}
							idMal
							id
							meanScore
							title {
								romaji
							}
							seasonYear
							type
						}
					}
				}
				seasonYear
			}
		}`, malID),
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var aniListResp models.AniListMangaResponse
	if err := json.NewDecoder(resp.Body).Decode(&aniListResp); err != nil {
		return nil, err
	}
	return &aniListResp, nil
}

func fetchAniListMangaRecommendations(malID int) (*models.AniListRecommendationResponse, error) {
	query := `query ($idMal: Int, $type: MediaType, $sort: [RecommendationSort]) {
		Media(idMal: $idMal, type: $type) {
			recommendations(sort: $sort) {
				nodes {
					mediaRecommendation {
						id
						idMal
						type
						title {
							romaji
							english
							native
						}
						coverImage {
							extraLarge
						}
					}
				}
			}
		}
	}`
	variables := map[string]interface{}{
		"idMal": malID,
		"type":  "MANGA",
		"sort":  []string{"RATING_DESC"},
	}
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[HATA]: Anilist API durum kodu: %d", resp.StatusCode)
	}
	var aniListResp models.AniListRecommendationResponse
	if err := json.NewDecoder(resp.Body).Decode(&aniListResp); err != nil {
		return nil, err
	}
	return &aniListResp, nil
}

func fetchAniListLightNovelData(malID int) (*models.AniListLightNovelResponse, error) {
	jsonData := map[string]interface{}{
		"query": fmt.Sprintf(`{
			Media(idMal: %d, type: MANGA, format: NOVEL) {
				id
				title {
					romaji
					english
					native
				}
				trending
				averageScore
				bannerImage
				coverImage {
					large
				}
				description
				tags {
					name
				}
				relations {
					edges {
						id
						relationType
						node {
							coverImage {
								extraLarge
							}
							idMal
							id
							meanScore
							title {
								romaji
							}
							seasonYear
							type
						}
					}
				}
				seasonYear
			}
		}`, malID),
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var aniListResp models.AniListLightNovelResponse
	if err := json.NewDecoder(resp.Body).Decode(&aniListResp); err != nil {
		return nil, err
	}
	return &aniListResp, nil
}
