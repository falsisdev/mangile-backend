package models

import "encoding/json"

type PageAsset struct {
	URL string `json:"url"`
}

type MangaPage struct {
	Asset PageAsset `json:"asset"`
}

type ChapterSource struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type MangaChapter struct {
	ChapterNumber float64        `json:"chapterNumber"`
	Title         string         `json:"title"`
	Source        *ChapterSource `json:"source,omitempty"`
	Pages         []MangaPage    `json:"pages"`
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

type AniListRelationNodeTitle struct {
	Romaji string `json:"romaji"`
}

type AniListRelationNodeCover struct {
	ExtraLarge string `json:"extraLarge"`
}

type AniListRelationNode struct {
	ID         int                      `json:"id"`
	IDMal      *int                     `json:"idMal"`
	Type       string                   `json:"type"`
	MeanScore  *float64                 `json:"meanScore"`
	SeasonYear *int                     `json:"seasonYear"`
	Title      AniListRelationNodeTitle `json:"title"`
	CoverImage AniListRelationNodeCover `json:"coverImage"`
}

type AniListRelationEdge struct {
	ID           int                 `json:"id"`
	RelationType string              `json:"relationType"`
	Node         AniListRelationNode `json:"node"`
}

type AniListRelationsConnection struct {
	Edges []AniListRelationEdge `json:"edges"`
}

func (c *AniListRelationsConnection) UnmarshalJSON(data []byte) error {
	type Alias AniListRelationsConnection
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var filtered []AniListRelationEdge
	for _, edge := range aux.Edges {
		if edge.Node.Type == "MANGA" {
			filtered = append(filtered, edge)
		}
	}
	c.Edges = filtered
	return nil
}

type AniListMangaMedia struct {
	AnilistID          int                        `json:"id"`
	AnilistTitle       AniListMangaTitle          `json:"title"`
	AnilistTrending    int                        `json:"trending"`
	AnilistScore       float64                    `json:"averageScore"`
	AniListBanner      string                     `json:"bannerImage"`
	AnilistCover       AniListMangaCoverImage     `json:"coverImage"`
	AnilistDescription string                     `json:"description"`
	AnilistTags        []interface{}              `json:"tags"`
	Relations          AniListRelationsConnection `json:"relations"`
	SeasonYear         int                        `json:"seasonYear"`
}

type AniListRecommendationResponse struct {
	Data AniListRecommendationData `json:"data"`
}

type AniListRecommendationData struct {
	Media AniListRecommendationMedia `json:"Media"`
}

type AniListRecommendationMedia struct {
	Recommendations AniListRecommendationConnection `json:"recommendations"`
}

type AniListRecommendationConnection struct {
	Nodes []AniListRecommendationNode `json:"nodes"`
}

type AniListRecommendationNode struct {
	MediaRecommendation AniListRecommendation `json:"mediaRecommendation"`
}

type AniListRecommendation struct {
	ID         int                             `json:"id"`
	IDMal      *int                            `json:"idMal"`
	Type       string                          `json:"type"`
	Title      AniListRecommendationTitle      `json:"title"`
	CoverImage AniListRecommendationCoverImage `json:"coverImage"`
}

type AniListRecommendationTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type AniListRecommendationCoverImage struct {
	ExtraLarge string `json:"extraLarge"`
}

type Manga struct {
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
	Chapters           []MangaChapter             `json:"chapters"`
	Notes              []interface{}              `json:"notes"`
}
