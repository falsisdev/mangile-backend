package models

type PageAsset struct {
	URL string `json:"url"`
}

type MangaPage struct {
	Asset PageAsset `json:"asset"`
}

type MangaChapter struct {
	ChapterNumber float64     `json:"chapterNumber"`
	Title         string      `json:"title"`
	Pages         []MangaPage `json:"pages"`
}

type SanityManga struct {
	ID                string         `json:"_id"`
	Type              string         `json:"_type"`
	CreatedAt         string         `json:"_createdAt"`
	UpdatedAt         string         `json:"_updatedAt"`
	MalID             int            `json:"myAnimeListId"`
	SanityBanner      string         `json:"bannerImage"`
	SanityCover       string         `json:"coverImage"`
	Chapters          []MangaChapter `json:"chapters"`
	Notes             []interface{}  `json:"notes"`
	SanityTitle       string         `json:"title"`
	SanityDescription string         `json:"description"`
	SanityTags        []string       `json:"tags"`
}

// JikanMangaResponse: Jikan veriyi genellikle bir "data" objesi içinde sarar
type JikanMangaResponse struct {
	Data JikanMangaData `json:"data"`
}

type JikanMangaData struct {
	MalURL           string        `json:"url"`
	MalTitleJapanese string        `json:"title_japanese"`
	MalTitleEnglish  string        `json:"title_english"`
	MalStatus        string        `json:"status"`
	MalScore         float64       `json:"score"`
	MalAuthors       []interface{} `json:"authors"`
	MalGenres        []interface{} `json:"genres"`
	MalThemes        []interface{} `json:"themes"`
}

// AniList GraphQL API kullandığı için dönen yanıt genellikle "data" -> "Media" şeklinde olur.
type AniListMangaResponse struct {
	Data AniListMangaData `json:"data"`
}

type AniListMangaData struct {
	Media AniListMangaMedia `json:"Media"`
}

type AniListMangaCoverImage struct {
	ExtraLarge string `json:"extraLarge"`
	Large      string `json:"large"`
	Medium     string `json:"medium"`
	Color      string `json:"color"`
}

type AniListMangaTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type AniListMangaMedia struct {
	AnilistID          int                    `json:"id"`
	AnilistTitle       AniListMangaTitle      `json:"title"`
	AnilistTrending    int                    `json:"trending"`
	AnilistScore       float64                `json:"averageScore"`
	AniListBanner      string                 `json:"bannerImage"`
	AnilistCover       AniListMangaCoverImage `json:"coverImage"`
	AnilistDescription string                 `json:"description"`
	AnilistTags        []interface{}          `json:"tags"`
}

type Manga struct {
	ID                 string         `json:"id"`
	Type               string         `json:"type"`
	SanityTitle        string         `json:"sanity_title"`
	SanityDescription  string         `json:"sanity_description"`
	SanityBanner       string         `json:"sanity_banner"`
	SanityCover        string         `json:"sanity_cover"`
	SanityTags         []string       `json:"sanity_tags"`
	AniListID          int            `json:"anilist_id"`
	AnilistTitle       string         `json:"anilist_title"`
	AnilistScore       float64        `json:"anilist_score"`
	AnilistBanner      string         `json:"anilist_banner"`
	AnilistCover       string         `json:"anilist_cover"`
	AnilistTags        []interface{}  `json:"anilist_tags"`
	AnilistDescription string         `json:"anilist_description"`
	AnilistTrending    int            `json:"anilist_trending"`
	MalID              int            `json:"mal_id"`
	MalTitleJapanese   string         `json:"title_japanese"`
	MalTitleEnglish    string         `json:"title_english"`
	MalStatus          string         `json:"mal_status"`
	MalScore           float64        `json:"mal_score"`
	MalAuthors         []interface{}  `json:"mal_authors"`
	MalGenres          []interface{}  `json:"mal_genres"`
	MalThemes          []interface{}  `json:"mal_themes"`
	MalURL             string         `json:"mal_url"`
	Chapters           []MangaChapter `json:"chapters"`
	Notes              []interface{}  `json:"notes"`
}
