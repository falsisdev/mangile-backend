package models

type Asset struct {
	URL string `json:"url"`
}

type Image struct {
	Asset Asset `json:"asset"`
}

type MangaChapter struct {
	ChapterNumber int     `json:"chapterNumber"`
	Title         string  `json:"title"`
	Pages         []Image `json:"pages"`
}

type Manga struct {
	ID          string         `json:"_id"`        //örn: 636cc2ad-6e8a-42d2-8bd8-10eddf4a2fa0
	Type        string         `json:"_type"`      //manga
	CreatedAt   string         `json:"_createdAt"` //örn: 2025-01-28T13:29:11Z
	UpdatedAt   string         `json:"_updatedAt"` //örn: 2025-07-22T13:58:40Z
	BannerImage string         `json:"bannerImage"`
	CoverImage  string         `json:"coverImage"`
	Chapters    []MangaChapter `json:"chapters"`
	Notes       []interface{}  `json:"notes"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	MALID       int            `json:"myAnimeListId"`
	Tags        []string       `json:"tags"`
}
