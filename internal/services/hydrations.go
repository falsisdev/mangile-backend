package services

import "github.com/falsisdev/mangile-backend/internal/models"

func hasJikanMangaData(data *models.JikanMangaData) bool {
	if data == nil {
		return false
	}
	return data.MalURL != "" || data.MalTitleJapanese != "" || data.MalTitleEnglish != "" || data.MalStatus != "" || data.MalScore > 0 || len(data.MalAuthors) > 0 || len(data.MalGenres) > 0 || len(data.MalThemes) > 0
}

func hasAniListMangaData(data *models.AniListMangaMedia) bool {
	if data == nil {
		return false
	}
	return data.AnilistID != 0 || data.AnilistTitle.Romaji != "" || data.AnilistTitle.English != "" || data.AnilistTitle.Native != "" || data.AnilistDescription != "" || data.AniListBanner != "" || data.AnilistCover.Large != "" || len(data.AnilistTags) > 0 || data.AnilistTrending != 0 || data.SeasonYear != 0 || len(data.Relations.Edges) > 0
}

func hasJikanLightNovelData(data *models.JikanLightNovelData) bool {
	if data == nil {
		return false
	}
	return data.MalURL != "" || data.MalTitleJapanese != "" || data.MalTitleEnglish != "" || data.MalStatus != "" || data.MalScore > 0 || len(data.MalAuthors) > 0 || len(data.MalGenres) > 0 || len(data.MalThemes) > 0
}

func hasAniListLightNovelData(data *models.AniListLightNovelMedia) bool {
	if data == nil {
		return false
	}
	return data.AnilistID != 0 || data.AnilistTitle.Romaji != "" || data.AnilistTitle.English != "" || data.AnilistTitle.Native != "" || data.AnilistDescription != "" || data.AniListBanner != "" || data.AnilistCover.Large != "" || len(data.AnilistTags) > 0 || data.AnilistTrending != 0 || data.SeasonYear != 0 || len(data.Relations.Edges) > 0
}

func hydrateMangaWithExternalData(manga *models.Manga, malID int) {
	if malID == 0 {
		return
	}

	manga.MalDataMissing = true
	manga.AniListDataMissing = true

	if jikanData, err := fetchJikanMangaData(malID); err == nil && jikanData != nil {
		if hasJikanMangaData(&jikanData.Data) {
			manga.MalDataMissing = false
			manga.MalURL = jikanData.Data.MalURL
			manga.MalTitleJapanese = jikanData.Data.MalTitleJapanese
			manga.MalTitleEnglish = jikanData.Data.MalTitleEnglish
			manga.MalStatus = jikanData.Data.MalStatus
			manga.MalScore = jikanData.Data.MalScore
			manga.MalAuthors = jikanData.Data.MalAuthors
			manga.MalGenres = jikanData.Data.MalGenres
			manga.MalThemes = jikanData.Data.MalThemes
		}
	}
	if aniListData, err := fetchAniListMangaData(malID); err == nil && aniListData != nil {
		media := aniListData.Data.Media
		if hasAniListMangaData(&media) {
			manga.AniListDataMissing = false
			manga.AniListID = media.AnilistID
			manga.AnilistTitle = media.AnilistTitle.Romaji
			if manga.AnilistTitle == "" {
				manga.AnilistTitle = media.AnilistTitle.English
			}
			manga.AnilistScore = media.AnilistScore
			manga.AnilistDescription = media.AnilistDescription
			manga.AnilistBanner = media.AniListBanner
			manga.AnilistCover = media.AnilistCover.Large
			manga.AnilistTags = media.AnilistTags
			manga.AnilistTrending = media.AnilistTrending
			manga.AnilistSeasonYear = media.SeasonYear
			manga.AnilistRelations = media.Relations
		}
	}
}

func hydrateLightNovelWithExternalData(lightNovel *models.LightNovel, malID int) {
	if malID == 0 {
		return
	}

	lightNovel.MalDataMissing = true
	lightNovel.AniListDataMissing = true

	if jikanData, err := fetchJikanMangaData(malID); err == nil && jikanData != nil {
		if hasJikanMangaData(&jikanData.Data) {
			lightNovel.MalDataMissing = false
			lightNovel.MalURL = jikanData.Data.MalURL
			lightNovel.MalTitleJapanese = jikanData.Data.MalTitleJapanese
			lightNovel.MalTitleEnglish = jikanData.Data.MalTitleEnglish
			lightNovel.MalStatus = jikanData.Data.MalStatus
			lightNovel.MalScore = jikanData.Data.MalScore
			lightNovel.MalAuthors = jikanData.Data.MalAuthors
			lightNovel.MalGenres = jikanData.Data.MalGenres
			lightNovel.MalThemes = jikanData.Data.MalThemes
		}
	}
	if aniListData, err := fetchAniListLightNovelData(malID); err == nil && aniListData != nil {
		media := aniListData.Data.Media
		if hasAniListLightNovelData(&media) {
			lightNovel.AniListDataMissing = false
			lightNovel.AniListID = media.AnilistID
			lightNovel.AnilistTitle = media.AnilistTitle.Romaji
			if lightNovel.AnilistTitle == "" {
				lightNovel.AnilistTitle = media.AnilistTitle.English
			}
			lightNovel.AnilistScore = media.AnilistScore
			lightNovel.AnilistDescription = media.AnilistDescription
			lightNovel.AnilistBanner = media.AniListBanner
			lightNovel.AnilistCover = media.AnilistCover.Large
			lightNovel.AnilistTags = media.AnilistTags
			lightNovel.AnilistTrending = media.AnilistTrending
			lightNovel.AnilistSeasonYear = media.SeasonYear
			lightNovel.AnilistRelations = media.Relations
		}
	}
}
