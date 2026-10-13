package models

type ReadingHistoryItem struct {
	Key           string  `json:"_key,omitempty"`
	SeriesID      string  `json:"seriesId"`
	SeriesTitle   string  `json:"seriesTitle"`
	SeriesType    string  `json:"seriesType"` // "manga" or "lightNovel"
	SeriesFormat  string  `json:"seriesFormat,omitempty"`
	CoverURL      string  `json:"coverUrl,omitempty"`
	ChapterID     string  `json:"chapterId"`
	ChapterNumber float64 `json:"chapterNumber"`
	VolumeNumber  float64 `json:"volumeNumber,omitempty"`
	ChapterTitle  string  `json:"chapterTitle,omitempty"`
	CurrentPage   int     `json:"currentPage,omitempty"`
	TotalPages    int     `json:"totalPages,omitempty"`
	Progress      float64 `json:"progress,omitempty"` // 0 to 100
	ReadAt        string  `json:"readAt"`
}

type ReadingHistoryUser struct {
	ID       string `json:"id"`
	LogtoID  string `json:"logtoId"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type ReadingHistoryDocument struct {
	ID        string               `json:"_id"`
	Type      string               `json:"_type"`
	UpdatedAt string               `json:"updatedAt,omitempty"`
	User      *ReadingHistoryUser  `json:"user,omitempty"`
	Items     []ReadingHistoryItem `json:"items"`
}

type UpsertHistoryRequest struct {
	LogtoID string             `json:"logtoId"`
	Item    ReadingHistoryItem `json:"item"`
}

type BatchSyncHistoryRequest struct {
	LogtoID string               `json:"logtoId"`
	Items   []ReadingHistoryItem `json:"items"`
}
