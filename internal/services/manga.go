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
		"bannerImage": bannerImage.asset->url,
		"coverImage": coverImage.asset->url,
		"chapters": chapters[]{
			chapterNumber,
			title,
			_key,
			"source": source -> {
				"id": _id,
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
		return nil, fmt.Errorf("İçerik bulunamadı veya ID eşleşmedi: %s", id)
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
			SanityTitle:       sanityData.SanityTitle,
			SanityDescription: sanityData.SanityDescription,
			SanityBanner:      sanityData.SanityBanner,
			SanityCover:       sanityData.SanityCover,
			SanityTags:        sanityData.SanityTags,
			MalID:             sanityData.MalID,
			Chapters:          sanityData.Chapters,
			Notes:             sanityData.Notes,
		}
		if finalLightNovel.MalID != 0 {
			if jikanData, err := fetchJikanMangaData(finalLightNovel.MalID); err == nil && jikanData != nil {
				finalLightNovel.MalURL = jikanData.Data.MalURL
				finalLightNovel.MalTitleJapanese = jikanData.Data.MalTitleJapanese
				finalLightNovel.MalTitleEnglish = jikanData.Data.MalTitleEnglish
				finalLightNovel.MalStatus = jikanData.Data.MalStatus
				finalLightNovel.MalScore = jikanData.Data.MalScore
				finalLightNovel.MalAuthors = jikanData.Data.MalAuthors
				finalLightNovel.MalGenres = jikanData.Data.MalGenres
				finalLightNovel.MalThemes = jikanData.Data.MalThemes
			}
			if aniListData, err := fetchAniListLightNovelData(finalLightNovel.MalID); err == nil && aniListData != nil {
				media := aniListData.Data.Media
				finalLightNovel.AniListID = media.AnilistID
				finalLightNovel.AnilistTitle = media.AnilistTitle.Romaji
				if finalLightNovel.AnilistTitle == "" {
					finalLightNovel.AnilistTitle = media.AnilistTitle.English
				}
				finalLightNovel.AnilistScore = media.AnilistScore
				finalLightNovel.AnilistDescription = media.AnilistDescription
				finalLightNovel.AnilistBanner = media.AniListBanner
				finalLightNovel.AnilistCover = media.AnilistCover.Large
				finalLightNovel.AnilistTags = media.AnilistTags
				finalLightNovel.AnilistTrending = media.AnilistTrending
				finalLightNovel.AnilistSeasonYear = media.SeasonYear
				finalLightNovel.AnilistRelations = media.Relations
			}
		}
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
	}
	if finalManga.MalID != 0 {
		if jikanData, err := fetchJikanMangaData(finalManga.MalID); err == nil && jikanData != nil {
			finalManga.MalURL = jikanData.Data.MalURL
			finalManga.MalTitleJapanese = jikanData.Data.MalTitleJapanese
			finalManga.MalTitleEnglish = jikanData.Data.MalTitleEnglish
			finalManga.MalStatus = jikanData.Data.MalStatus
			finalManga.MalScore = jikanData.Data.MalScore
			finalManga.MalAuthors = jikanData.Data.MalAuthors
			finalManga.MalGenres = jikanData.Data.MalGenres
			finalManga.MalThemes = jikanData.Data.MalThemes
		}
		if aniListData, err := fetchAniListMangaData(finalManga.MalID); err == nil && aniListData != nil {
			media := aniListData.Data.Media
			finalManga.AniListID = media.AnilistID
			finalManga.AnilistTitle = media.AnilistTitle.Romaji
			if finalManga.AnilistTitle == "" {
				finalManga.AnilistTitle = media.AnilistTitle.English
			}
			finalManga.AnilistScore = media.AnilistScore
			finalManga.AnilistDescription = media.AnilistDescription
			finalManga.AnilistBanner = media.AniListBanner
			finalManga.AnilistCover = media.AnilistCover.Large
			finalManga.AnilistTags = media.AnilistTags
			finalManga.AnilistTrending = media.AnilistTrending
			finalManga.AnilistSeasonYear = media.SeasonYear
			finalManga.AnilistRelations = media.Relations
		}
	}
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