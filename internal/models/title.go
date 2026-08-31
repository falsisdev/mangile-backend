package models

import (
	"encoding/json"
	"time"
)

// ==========================================
// SANITY
// ==========================================

type ImageAsset struct {
	URL string `json:"url,omitempty"`
}

type Source struct {
	ID         string     `json:"_id,omitempty"`
	Type       string     `json:"_type,omitempty"`
	Name       string     `json:"name,omitempty"`
	Website    string     `json:"website,omitempty"`
	CoverImage ImageAsset `json:"coverImage,omitempty"`
	Logo       ImageAsset `json:"logo,omitempty"`
}

type Chapter struct {
	ID            string       `json:"_id,omitempty"`
	CreatedAt     *time.Time   `json:"_createdAt,omitempty"`
	Title         string       `json:"title,omitempty"`
	VolumeNumber  float64      `json:"volumeNumber"`
	ChapterNumber float64      `json:"chapterNumber"`
	Pages         []ImageAsset `json:"pages,omitempty"`
	Source        *Source      `json:"source,omitempty"`
}

type TitleMedia struct {
	ID                string       `json:"_id,omitempty"`
	Type              string       `json:"_type,omitempty"`
	CreatedAt         time.Time    `json:"_createdAt"`
	UpdatedAt         time.Time    `json:"_updatedAt"`
	Title             string       `json:"title,omitempty"`
	Slug              string       `json:"slug,omitempty"`
	Description       string       `json:"description,omitempty"`
	MyAnimeListID     int          `json:"myAnimeListId,omitempty"`
	UploadStatus      string       `json:"uploadStatus,omitempty"`
	Notes             []any        `json:"notes,omitempty"`
	Tags              []string     `json:"tags,omitempty"`
	Format            string       `json:"format,omitempty"`
	BannerImage       ImageAsset   `json:"bannerImage"`
	CoverImage        ImageAsset   `json:"coverImage"`
	Chapters          []Chapter    `json:"chapters,omitempty"`
	ExternalMAL       *JikanData   `json:"externalMal,omitempty"`
	ExternalAniList   *AniListData `json:"externalAnilist,omitempty"`
	SanityDataMissing bool         `json:"sanity_data_missing"`
}

// ==========================================
// JIKAN
// ==========================================

type JikanResponse struct {
	Data JikanData `json:"data"`
}

type JikanData struct {
	MalURL           string  `json:"url"`
	MalType          string  `json:"type"`
	MalTitleJapanese string  `json:"title_japanese"`
	MalTitleEnglish  string  `json:"title_english"`
	MalStatus        string  `json:"status"`
	MalScore         float64 `json:"score"`
	MalAuthors       []any   `json:"authors"`
	MalGenres        []any   `json:"genres"`
	MalThemes        []any   `json:"themes"`
}

// ==========================================
// ANILIST
// ==========================================

type AniListResponse struct {
	Data AniListResultData `json:"data"`
}

type AniListResultData struct {
	Media AniListData `json:"Media"`
}

type AniListCoverImage struct {
	ExtraLarge string `json:"extraLarge"`
	Large      string `json:"large"`
	Medium     string `json:"medium"`
	Color      string `json:"color"`
}

type AniListTitle struct {
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

type AniListData struct {
	AnilistID          int                              `json:"id"`
	CountryOfOrigin    string                           `json:"countryOfOrigin,omitempty"`
	Format             string                           `json:"format,omitempty"`
	AnilistTitle       AniListTitle                     `json:"title"`
	AnilistTrending    int                              `json:"trending"`
	AnilistScore       float64                          `json:"averageScore"`
	AniListBanner      string                           `json:"bannerImage"`
	AnilistCover       AniListCoverImage                `json:"coverImage"`
	AnilistDescription string                           `json:"description"`
	AnilistTags        []any                            `json:"tags"`
	Relations          AniListRelationsConnection       `json:"relations"`
	Recommendations    *AniListRecommendationConnection `json:"recommendations,omitempty"`
	SeasonYear         int                              `json:"seasonYear"`
}

// ==========================================
// ANILIST RECOMMENDATION
// ==========================================

type AniListRecommendationResponse struct {
	Data AniListRecommendationData `json:"data"`
}

type AniListRecommendationData struct {
	Media struct {
		Recommendations AniListRecommendationConnection `json:"recommendations"`
	} `json:"Media"`
}

type AniListRecommendationConnection struct {
	Nodes []AniListRecommendationNode `json:"nodes"`
}

type AniListRecommendationNode struct {
	MediaRecommendation AniListRecommendation `json:"mediaRecommendation"`
}

type AniListRecommendation struct {
	ID         int                      `json:"id"`
	IDMal      *int                     `json:"idMal"`
	Type       string                   `json:"type"`
	Title      AniListTitle             `json:"title"`
	CoverImage AniListRelationNodeCover `json:"coverImage"`
}
