package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func GetMangaList(filterType string, limit int, page int, searchQuery string) ([]models.MangaCard, error) {
	if page < 1 {
		page = 1
	}

	var sortParam string
	var statusParam *string

	if searchQuery != "" {
		sortParam = ""
	} else {
		switch filterType {
		case "POPULAR":
			sortParam = "sort: [POPULARITY_DESC],"
		case "HIGHEST_SCORE":
			sortParam = "sort: [SCORE_DESC],"
		case "TRENDING":
			sortParam = "sort: [TRENDING_DESC],"
		case "UPCOMING":
			sortParam = "sort: [START_DATE_DESC],"
			statusVal := "NOT_YET_RELEASED"
			statusParam = &statusVal
		default:
			sortParam = "sort: [POPULARITY_DESC],"
		}
	}

	query := fmt.Sprintf(`
	query Media($type: MediaType, $isAdult: Boolean, $countryOfOrigin: CountryCode, $page: Int, $perPage: Int, $status: MediaStatus, $search: String) {
		Page (page: $page, perPage: $perPage) {
			media (type: $type, %s search: $search, isAdult: $isAdult, countryOfOrigin: $countryOfOrigin, status: $status) {
				id
				idMal
				type
				format
				status
				meanScore
				bannerImage
				description
				startDate {
					year
				}
				coverImage {
					large
				}
				title {
					romaji
					english
					native
				}
			}
		}
	}`, sortParam)

	variables := map[string]interface{}{
		"type":            "MANGA",
		"page":            page,
		"perPage":         limit,
		"isAdult":         false,
		"countryOfOrigin": "JP",
		"status":          statusParam,
	}

	if statusParam == nil {
		delete(variables, "status")
	}

	if searchQuery != "" {
		variables["search"] = searchQuery
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, fmt.Errorf("[HATA]: Request body marshalling failed: %w", err)
	}

	req, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("[HATA]: HTTP isteği oluşturulurken bir sorun oluştu: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[HATA]: Anilist API isteğinde bir sorun oluştu: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("[HATA]: Anilist API durum kodu: %d\nyanıt: %s", resp.StatusCode, string(bodyBytes))
	}

	var aniListResp models.AniListListResponse
	if err := json.NewDecoder(resp.Body).Decode(&aniListResp); err != nil {
		return nil, fmt.Errorf("[HATA]: Anilist yanıtı çözümlenirken hata oluştu: %w", err)
	}

	var malIDs []string
	for _, media := range aniListResp.Data.Page.Media {
		if media.IDMal != 0 {
			malIDs = append(malIDs, strconv.Itoa(media.IDMal))
		}
	}

	sanityMatches := make(map[int]struct {
		Description string
		BannerImage string
	})

	projectID := os.Getenv("SANITY_PROJECT_ID")
	if len(malIDs) > 0 && projectID != "" {
		idListStr := "[" + strings.Join(malIDs, ",") + "]"

		sanityQuery := fmt.Sprintf(`*[_type == "manga" && myAnimeListId in %s]{
			myAnimeListId,
			description,
			"bannerImage": bannerImage.asset->url
		}`, idListStr)

		baseURL := fmt.Sprintf("https://%s.api.sanity.io/v2021-10-21/data/query/production", projectID)
		u, _ := url.Parse(baseURL)
		q := u.Query()
		q.Set("query", sanityQuery)
		u.RawQuery = q.Encode()

		if sResp, sErr := http.Get(u.String()); sErr == nil {
			defer sResp.Body.Close()
			if sResp.StatusCode == http.StatusOK {
				var sanityListWrapper struct {
					Result []struct {
						MyAnimeListId int    `json:"myAnimeListId"`
						Description   string `json:"description"`
						BannerImage   string `json:"bannerImage"`
					} `json:"result"`
				}
				if json.NewDecoder(sResp.Body).Decode(&sanityListWrapper) == nil {
					for _, item := range sanityListWrapper.Result {
						sanityMatches[item.MyAnimeListId] = struct {
							Description string
							BannerImage string
						}{
							Description: item.Description,
							BannerImage: item.BannerImage,
						}
					}
				}
			}
		}
	}

	var mangaCards []models.MangaCard
	for _, media := range aniListResp.Data.Page.Media {
		mainTitle := media.Title.Romaji
		if mainTitle == "" {
			mainTitle = media.Title.English
		}

		bannerImg := media.BannerImage
		var sanityDesc string
		var hasLocal bool

		if localData, exists := sanityMatches[media.IDMal]; exists {
			hasLocal = true
			sanityDesc = localData.Description
			if localData.BannerImage != "" {
				bannerImg = localData.BannerImage
			}
		}

		card := models.MangaCard{
			AniListID:          media.ID,
			MyAnimeListID:      media.IDMal,
			AniListTitle:       mainTitle,
			TitleRomaji:        media.Title.Romaji,
			TitleEnglish:       media.Title.English,
			TitleNative:        media.Title.Native,
			Type:               media.Type,
			Format:             media.Format,
			Status:             media.Status,
			Score:              media.MeanScore,
			CoverImage:         media.CoverImage.Large,
			BannerImage:        bannerImg,
			AniListDescription: media.Description,
			MalType:            media.Format,
			MalYear:            media.StartDate.Year,
			HasLocalContent:    hasLocal,
			SanityDescription:  sanityDesc,
		}
		mangaCards = append(mangaCards, card)
	}
	return mangaCards, nil
}