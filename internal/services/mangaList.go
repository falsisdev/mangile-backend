package services

import (
	"context"
	"net/http"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

func GetMangaList(ctx context.Context, filterType string, limit int, page int, searchQuery string, sortParam string, formatParam string, genreParam string, statusParam string) ([]models.MangaCard, error) {
	if page < 1 {
		page = 1
	}

	var sort []string
	var status any

	if sortParam != "" {
		sort = []string{sortParam}
	} else if searchQuery == "" {
		switch filterType {
		case "POPULAR":
			sort = []string{"POPULARITY_DESC"}
		case "HIGHEST_SCORE":
			sort = []string{"SCORE_DESC"}
		case "TRENDING":
			sort = []string{"TRENDING_DESC"}
		case "UPCOMING":
			sort = []string{"START_DATE_DESC"}
			status = "NOT_YET_RELEASED"
		default:
			sort = []string{"POPULARITY_DESC"}
		}
	}

	if statusParam != "" {
		status = statusParam
	}

	variables := map[string]any{
		"page":    page,
		"perPage": limit,
		"search":  nil,
		"sort":    sort,
	}

	if status != nil {
		variables["status"] = status
	}
	if formatParam != "" {
		variables["format"] = formatParam
	}
	if genreParam != "" {
		variables["genre"] = genreParam
	}

	if searchQuery != "" {
		variables["search"] = searchQuery
	}

	var aniListResp models.AniListListResponse
	httpClient := &http.Client{Timeout: 15 * time.Second}
	if err := postAniListQuery(ctx, httpClient, queries.AniListMangaList, variables, &aniListResp); err != nil {
		return nil, err
	}

	malIDs := make([]int, 0, len(aniListResp.Data.Page.Media))
	for _, media := range aniListResp.Data.Page.Media {
		if media.IDMal != 0 {
			malIDs = append(malIDs, media.IDMal)
		}
	}

	sanityMatches := make(map[int]struct {
		Description string
		BannerImage string
	})

	// Sanity zenginlestirmesi zorunlu degil; erisilemezse AniList verisiyle devam edilir.
	if len(malIDs) > 0 {
		if sanity, sErr := newSanityClient(); sErr == nil {
			var localTitles []struct {
				MyAnimeListId int    `json:"myAnimeListId"`
				Description   string `json:"description"`
				BannerImage   string `json:"bannerImage"`
			}

			if qErr := sanity.query(ctx, queries.MangaListQuery, map[string]any{"ids": malIDs}, &localTitles); qErr == nil {
				for _, item := range localTitles {
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
