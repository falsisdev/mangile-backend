package models

// LatestTitle, /api/latestTitles yanıtındaki tek bir yerel seriyi temsil eder.
type LatestTitle struct {
	ID          string   `json:"_id"`
	Type        string   `json:"_type"`
	CreatedAt   string   `json:"_createdAt"`
	UpdatedAt   string   `json:"_updatedAt"`
	MalID       int      `json:"myAnimeListId"`
	BannerImage string   `json:"bannerImage"`
	CoverImage  string   `json:"coverImage"`
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
}
