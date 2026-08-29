package models

import "time"

// TitleByTag, /api/titlesByTag yanıtındaki tek bir eseri (manga/lightNovel)
// temsil eder.
type TitleByTag struct {
	ID            string     `json:"_id,omitempty"`
	Type          string     `json:"_type,omitempty"`
	CreatedAt     time.Time  `json:"_createdAt"`
	UpdatedAt     time.Time  `json:"_updatedAt"`
	Title         string     `json:"title,omitempty"`
	Description   string     `json:"description,omitempty"`
	Slug          string     `json:"slug,omitempty"`
	MyAnimeListID int        `json:"myAnimeListId,omitempty"`
	Tags          []string   `json:"tags,omitempty"`
	CoverImage    ImageAsset `json:"coverImage"`
}
