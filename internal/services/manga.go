package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func GetManga(id string) (interface{}, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("SANITY_PROJECT_ID ortam değişkeni bulunamadı")
	}

	query := fmt.Sprintf(`*[_type in ["manga", "lightNovel"] && myAnimeListId == %s][0]{
		_id,
		_type,
		_createdAt,
		_updatedAt,
		myAnimeListId,
		title,
		description,
		tags,
		uploadStatus,
		"bannerImage": bannerImage.asset->url,
		"coverImage": coverImage.asset->url,
		"chapters": chapters[]{
			chapterNumber,
			title,
			_key,
			"source": source -> {
				_id,
				name
			},
			"pages": pages[]{
				"asset": {
					"url": asset->url
				}
			},
			content
		},
		notes
	}`, id)

	baseURL := fmt.Sprintf("https://%s.api.sanity.io/v2021-10-21/data/query/production", projectID)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawWrapper struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawWrapper); err != nil {
		return nil, err
	}

	if len(rawWrapper.Result) == 0 || string(rawWrapper.Result) == "null" {
		malID, _ := strconv.Atoi(id)
		fallbackManga := buildMissingSanityManga(id, malID)
		hydrateMangaWithExternalData(fallbackManga, fallbackManga.MalID)
		return fallbackManga, nil
	}

	var typeCheck struct {
		Type string `json:"_type"`
	}
	if err := json.Unmarshal(rawWrapper.Result, &typeCheck); err != nil {
		return nil, err
	}

	if typeCheck.Type == "lightNovel" {
		var sanityData models.SanityLightNovel
		if err := json.Unmarshal(rawWrapper.Result, &sanityData); err != nil {
			return nil, err
		}
		finalLightNovel := &models.LightNovel{
			ID:                sanityData.ID,
			Type:              sanityData.Type,
			UploadStatus:      sanityData.UploadStatus,
			SanityTitle:       sanityData.SanityTitle,
			SanityDescription: sanityData.SanityDescription,
			SanityBanner:      sanityData.SanityBanner,
			SanityCover:       sanityData.SanityCover,
			SanityTags:        sanityData.SanityTags,
			MalID:             sanityData.MalID,
			Chapters:          sanityData.Chapters,
			Notes:             sanityData.Notes,
		}
		hydrateLightNovelWithExternalData(finalLightNovel, finalLightNovel.MalID)
		return finalLightNovel, nil
	}

	var sanityData models.SanityManga
	if err := json.Unmarshal(rawWrapper.Result, &sanityData); err != nil {
		return nil, err
	}
	finalManga := &models.Manga{
		ID:                sanityData.ID,
		Type:              sanityData.Type,
		SanityTitle:       sanityData.SanityTitle,
		SanityDescription: sanityData.SanityDescription,
		SanityBanner:      sanityData.SanityBanner,
		SanityCover:       sanityData.SanityCover,
		SanityTags:        sanityData.SanityTags,
		MalID:             sanityData.MalID,
		Chapters:          sanityData.Chapters,
		Notes:             sanityData.Notes,
		UploadStatus:      sanityData.UploadStatus,
	}
	hydrateMangaWithExternalData(finalManga, finalManga.MalID)
	return finalManga, nil
}

func GetMangaRecommendations(id string) ([]models.AniListRecommendation, error) {
	malID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("[HATA]: Geçerli bir MAL ID'si girilmedi: %w", err)
	}
	aniListResp, err := fetchAniListMangaRecommendations(malID)
	if err != nil {
		return nil, fmt.Errorf("[HATA]: Anilist önerileri alınırken bir hata oluştu: %w", err)
	}
	if aniListResp == nil {
		return nil, nil
	}
	var recommendations []models.AniListRecommendation
	for _, node := range aniListResp.Data.Media.Recommendations.Nodes {
		recommendations = append(recommendations, node.MediaRecommendation)
	}
	return recommendations, nil
}

func buildMissingSanityManga(id string, malID int) *models.Manga {
	return &models.Manga{
		ID:                id,
		Type:              "manga",
		SanityDataMissing: true,
		MalID:             malID,
	}
}

func buildMissingSanityLightNovel(id string, malID int) *models.LightNovel {
	return &models.LightNovel{
		ID:                id,
		Type:              "lightNovel",
		SanityDataMissing: true,
		MalID:             malID,
	}
}
