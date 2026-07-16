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

func GetSanityList(filterType string) ([]models.SanityList, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	query := fmt.Sprintf(`*[_type == 'manga' || _type == 'lightNovel'] | order(_%s desc){
		"id": _id,
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

func GetUser(id string) (*models.User, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	query := fmt.Sprintf(`*[_type == "auth" && logtoId == "%s"][0]{
	"id": _id,
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

func GetList(id string) (*models.List, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	query := fmt.Sprintf(`*[_type == "lists" && _id == "%s"][0]{
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
	"id": _id,
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

func GetManga(id string) (*models.Manga, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("SANITY_PROJECT_ID ortam değişkeni bulunamadı")
	}
	query := fmt.Sprintf(`*[_type == "manga" && myAnimeListId == %s][0]{
		"id": _id,
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
			}
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
	var sanityWrapper struct {
		Result models.SanityManga `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sanityWrapper); err != nil {
		return nil, err
	}
	sanityData := sanityWrapper.Result
	if sanityData.ID == "" {
		return nil, fmt.Errorf("Manga bulunamadı veya ID eşleşmedi: %s", id)
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

func GetLightNovel(id string) (*models.LightNovel, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("SANITY_PROJECT_ID ortam değişkeni bulunamadı")
	}
	query := fmt.Sprintf(`*[_type == "lightNovel" && myAnimeListId == %s][0]{
	"id": _id, 
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
		_key,
		"source": source -> {
			"id": _id,
			name
		},
		"content": content
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
	var sanityWrapper struct {
		Result models.SanityLightNovel `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sanityWrapper); err != nil {
		return nil, err
	}
	sanityData := sanityWrapper.Result
	if sanityData.ID == "" {
		return nil, fmt.Errorf("Light Novel bulunamadı: %s", id)
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
	var aniListResp models.AniListLightNovelResponse
	if err := json.NewDecoder(resp.Body).Decode(&aniListResp); err != nil {
		return nil, err
	}
	return &aniListResp, nil
}
