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

func GetMangaList(filterType string, limit int) ([]models.MangaCard, error) {
	var sortParam string
	var statusParam string

	switch filterType {
	case "POPULAR":
		sortParam = "POPULARITY_DESC"
	case "HIGHEST_SCORE":
		sortParam = "SCORE_DESC"
	case "TRENDING":
		sortParam = "TRENDING_DESC"
	case "UPCOMING":
		sortParam = "START_DATE_DESC"
		statusParam = "NOT_YET_RELEASED"
	default:
		sortParam = "POPULARITY_DESC"
	}

	query := `
	query ($page: Int, $perPage: Int, $sort: [MediaSort], $status: MediaStatus) {
		Page (page: $page, perPage: $perPage) {
			media (type: MANGA, sort: $sort, status: $status) {
				id
				idMal
				type
				format
				status
				meanScore
				bannerImage
				description
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
	}`

	variables := map[string]interface{}{
		"page":    1,
		"perPage": limit,
		"sort":    []string{sortParam},
	}
	if statusParam != "" {
		variables["status"] = statusParam
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, fmt.Errorf("request body marshalling failed: %w", err)
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
		return nil, fmt.Errorf("[HATA]: Anilist API durum kodu: %d", resp.StatusCode)
	}

	var aniListResp models.AniListListResponse
	if err := json.NewDecoder(resp.Body).Decode(&aniListResp); err != nil {
		return nil, fmt.Errorf("[HATA]: Anilist yanıtı çözümlenirken hata oluştu: %w", err)
	}

	var mangaCards []models.MangaCard

	for _, media := range aniListResp.Data.Page.Media {
		mainTitle := media.Title.Romaji
		if mainTitle == "" {
			mainTitle = media.Title.English
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
			BannerImage:        media.BannerImage,
			AniListDescription: media.Description,
		}
		mangaCards = append(mangaCards, card)
	}

	return mangaCards, nil
}

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

func GetManga(id string) (*models.Manga, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("SANITY_PROJECT_ID ortam değişkeni bulunamadı")
	}

	// 1. Sanity GROQ Sorgusu
	query := fmt.Sprintf(`*[_type == "manga" && myAnimeListId == %s][0]{
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
		return nil, fmt.Errorf("Manga bulunamadı: %s", id)
	}

	finalLightNovel := &models.Manga{
		ID:                sanityData.ID,
		Type:              sanityData.Type,
		SanityTitle:       sanityData.SanityTitle,
		SanityDescription: sanityData.SanityDescription,
		SanityBanner:      sanityData.SanityBanner,
		SanityCover:       sanityData.SanityCover,
		SanityTags:        sanityData.SanityTags,
		MalID:             sanityData.MalID,
		Chapters:          sanityData.Chapters, // Ortak paylaşılan şema
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

		if aniListData, err := fetchAniListMangaData(finalLightNovel.MalID); err == nil && aniListData != nil {
			media := aniListData.Data.Media
			finalLightNovel.AniListID = media.AnilistID
			finalLightNovel.AnilistBanner = media.AniListBanner
			finalLightNovel.AnilistCover = media.AnilistCover.Large
			finalLightNovel.AnilistTags = media.AnilistTags
			finalLightNovel.AnilistTrending = media.AnilistTrending
		}
	}

	return finalLightNovel, nil
}

func GetLightNovel(id string) (*models.LightNovel, error) {
	projectID := os.Getenv("SANITY_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("SANITY_PROJECT_ID ortam değişkeni bulunamadı")
	}

	// 1. Sanity GROQ Sorgusu
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
		Chapters:          sanityData.Chapters, // Ortak paylaşılan şema
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

		if aniListData, err := fetchAniListMangaData(finalLightNovel.MalID); err == nil && aniListData != nil {
			media := aniListData.Data.Media
			finalLightNovel.AniListID = media.AnilistID
			finalLightNovel.AnilistBanner = media.AniListBanner
			finalLightNovel.AnilistCover = media.AnilistCover.Large
			finalLightNovel.AnilistTags = media.AnilistTags
			finalLightNovel.AnilistTrending = media.AnilistTrending
		}
	}

	return finalLightNovel, nil
}

// -------------------------- FETCH MANGA --------------------------------------
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
	jsonData := map[string]string{
		"query": fmt.Sprintf(`{
			Media(idMal: %d, type: MANGA) {
				id
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

//-------------------------- FETCH MANGA --------------------------------------

// -------------------------- FETCH LIGHTNOVEL --------------------------------------
func fetchJikanLightNovelData(malID int) (*models.JikanLightNovelResponse, error) {
	url := fmt.Sprintf("https://api.jikan.moe/v4/manga/%d", malID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jikanResp models.JikanLightNovelResponse
	if err := json.NewDecoder(resp.Body).Decode(&jikanResp); err != nil {
		return nil, err
	}
	return &jikanResp, nil
}

func fetchAniListLightNovelData(malID int) (*models.AniListLightNovelResponse, error) {
	jsonData := map[string]string{
		"query": fmt.Sprintf(`{
			Media(idMal: %d, type: MANGA) {
				id
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

//-------------------------- FETCH LIGHTNOVEL --------------------------------------
