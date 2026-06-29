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

type Manga struct {
	ID            string         `json:"_id"`        //örn: 636cc2ad-6e8a-42d2-8bd8-10eddf4a2fa0
	Type          string         `json:"_type"`      //manga
	CreatedAt     string         `json:"_createdAt"` //örn: 2025-01-28T13:29:11Z
	UpdatedAt     string         `json:"_updatedAt"` //örn: 2025-07-22T13:58:40Z
	BannerImage   string         `json:"bannerImage"`
	CoverImage    string         `json:"coverImage"`
	Chapters      []MangaChapter `json:"chapters"`
	Notes         []interface{}  `json:"notes"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	MyAnimeListID int            `json:"myAnimeListId"`
	Tags          []string       `json:"tags"`
}
