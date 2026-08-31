package models

// LocalTitle, /api/localTitles yanıtındaki tek bir yerel seriyi temsil eder.
type LocalTitle struct {
	ID           string   `json:"_id"`
	Type         string   `json:"_type"`
	CreatedAt    string   `json:"_createdAt"`
	UpdatedAt    string   `json:"_updatedAt"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	MalID        int      `json:"myAnimeListId"`
	Description  string   `json:"description"`
	UploadStatus string   `json:"uploadStatus"`
	Tags         []string `json:"tags"`
	CoverImage   string   `json:"coverImage"`
	BannerImage  string   `json:"bannerImage"`
}

// LocalTitlesQueryResult, Sanity GROQ sorgusundan dönenham sonucu temsil eder.
type LocalTitlesQueryResult struct {
	Data  []LocalTitle `json:"data"`
	Total int          `json:"total"`
}

// LocalTitlesResponse, /api/localTitles HTTP yanıtının yapısıdır.
type LocalTitlesResponse struct {
	Data  []LocalTitle `json:"data"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}
