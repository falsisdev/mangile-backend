package models

import "time"

type ChapterTitleRef struct {
	ID            string   `json:"_id,omitempty"`
	Type          string   `json:"_type,omitempty"`
	MyAnimeListID int      `json:"myAnimeListId,omitempty"`
	Title         string   `json:"title,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Format        string      `json:"format,omitempty"`
	CoverImage    *ImageAsset `json:"coverImage,omitempty"`
}

type ChapterListItem struct {
	ID            string  `json:"_id,omitempty"`
	Title         string  `json:"title,omitempty"`
	ChapterNumber float64 `json:"chapterNumber"`
	VolumeNumber  float64 `json:"volumeNumber"`
	Source        *Source `json:"source,omitempty"`
}

type ChapterDetails struct {
	ID            string            `json:"_id,omitempty"`
	Type          string            `json:"_type,omitempty"`
	CreatedAt     *time.Time        `json:"_createdAt,omitempty"`
	UpdatedAt     *time.Time        `json:"_updatedAt,omitempty"`
	Title         string            `json:"title,omitempty"`
	ChapterNumber float64           `json:"chapterNumber"`
	VolumeNumber  float64           `json:"volumeNumber"`
	Pages         []ImageAsset      `json:"pages,omitempty"`
	Content       []any             `json:"content,omitempty"`
	Manga         *ChapterTitleRef  `json:"manga,omitempty"`
	LightNovel    *ChapterTitleRef  `json:"lightNovel,omitempty"`
	Source        *Source           `json:"source,omitempty"`
	Chapters      []ChapterListItem `json:"chapters,omitempty"`
}
