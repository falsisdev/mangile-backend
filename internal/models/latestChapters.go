package models

import "time"

type LatestChapterTitleRef struct {
	ID            string     `json:"_id,omitempty"`
	Title         string     `json:"title,omitempty"`
	MyAnimeListID int        `json:"myAnimeListId,omitempty"`
	CoverImage    ImageAsset `json:"coverImage"`
}

type LatestChapter struct {
	ID            string                 `json:"_id,omitempty"`
	Type          string                 `json:"_type,omitempty"`
	CreatedAt     *time.Time             `json:"_createdAt,omitempty"`
	UpdatedAt     *time.Time             `json:"_updatedAt,omitempty"`
	Title         string                 `json:"title,omitempty"`
	ChapterNumber float64                `json:"chapterNumber"`
	VolumeNumber  float64                `json:"volumeNumber"`
	Manga         *LatestChapterTitleRef `json:"manga,omitempty"`
	LightNovel    *LatestChapterTitleRef `json:"lightNovel,omitempty"`
	Source        *Source                `json:"source,omitempty"`
}
