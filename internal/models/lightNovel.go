package models

type LightNovelChapter struct {
	ChapterNumber float64       `json:"chapterNumber"`
	Title         string        `json:"title"`
	Content       []interface{} `json:"content"`
}

type SanityLightNovel struct {
	ID                string              `json:"_id"`
	Type              string              `json:"_type"`
	CreatedAt         string              `json:"_createdAt"`
	UpdatedAt         string              `json:"_updatedAt"`
	SanityBanner      string              `json:"bannerImage"`
	SanityCover       string              `json:"coverImage"`
	Chapters          []LightNovelChapter `json:"chapters"`
	Notes             []interface{}       `json:"notes"`
	SanityTitle       string              `json:"title"`
	SanityDescription string              `json:"description"`
	MalID             int                 `json:"myAnimeListId"`
	SanityTags        []string            `json:"tags"`
}

type JikanLightNovelResponse struct {
	Data JikanLightNovelData `json:"data"`
}

type JikanLightNovelData struct {
	MalURL           string        `json:"url"`
	MalTitleJapanese string        `json:"title_japanese"`
	MalTitleEnglish  string        `json:"title_english"`
	MalStatus        string        `json:"status"`
	MalScore         float64       `json:"score"`
	MalAuthors       []interface{} `json:"authors"`
	MalGenres        []interface{} `json:"genres"`
	MalThemes        []interface{} `json:"themes"`
}

type AniListLightNovelResponse struct {
	Data AniListLightNovelData `json:"data"`
}

type AniListLightNovelData struct {
	Media AniListLightNovelMedia `json:"Media"`
}

type AniListLightNovelCoverImage struct {
	ExtraLarge string `json:"extraLarge"`
	Large      string `json:"large"`
	Medium     string `json:"medium"`
	Color      string `json:"color"`
}

type AniListLightNovelTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type AniListLightNovelMedia struct {
	AnilistID          int                         `json:"id"`
	AnilistTitle       AniListLightNovelTitle      `json:"title"`
	AnilistTrending    int                         `json:"trending"`
	AnilistScore       float64                     `json:"averageScore"`
	AniListBanner      string                      `json:"bannerImage"`
	AnilistCover       AniListLightNovelCoverImage `json:"coverImage"`
	AnilistDescription string                      `json:"description"`
	AnilistTags        []interface{}               `json:"tags"`
	Relations          AniListRelationsConnection  `json:"relations"`
	SeasonYear         int                         `json:"seasonYear"`
}

type LightNovel struct {
	ID                 string                     `json:"id"`
	Type               string                     `json:"type"`
	SanityTitle        string                     `json:"sanity_title"`
	SanityDescription  string                     `json:"sanity_description"`
	SanityBanner       string                     `json:"sanity_banner"`
	SanityCover        string                     `json:"sanity_cover"`
	SanityTags         []string                   `json:"sanity_tags"`
	AniListID          int                        `json:"anilist_id"`
	AnilistTitle       string                     `json:"anilist_title"`
	AnilistScore       float64                    `json:"anilist_score"`
	AnilistBanner      string                     `json:"anilist_banner"`
	AnilistCover       string                     `json:"anilist_cover"`
	AnilistTags        []interface{}              `json:"anilist_tags"`
	AnilistDescription string                     `json:"anilist_description"`
	AnilistTrending    int                        `json:"anilist_trending"`
	AnilistSeasonYear  int                        `json:"anilist_season_year"`
	AnilistRelations   AniListRelationsConnection `json:"anilist_relations"`
	MalID              int                        `json:"mal_id"`
	MalTitleJapanese   string                     `json:"title_japanese"`
	MalTitleEnglish    string                     `json:"title_english"`
	MalStatus          string                     `json:"mal_status"`
	MalScore           float64                    `json:"mal_score"`
	MalAuthors         []interface{}              `json:"mal_authors"`
	MalGenres          []interface{}              `json:"mal_genres"`
	MalThemes          []interface{}              `json:"mal_themes"`
	MalURL             string                     `json:"mal_url"`
	Chapters           []LightNovelChapter        `json:"chapters"`
	Notes              []interface{}              `json:"notes"`
}
