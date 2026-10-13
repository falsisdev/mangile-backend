package models

type SeriesTrackerItem struct {
	Key          string  `json:"_key,omitempty"`
	SeriesID     string  `json:"seriesId"`
	SeriesTitle  string  `json:"seriesTitle,omitempty"`
	MALID        int     `json:"malId,omitempty"`
	AniListID    int     `json:"anilistId,omitempty"`
	MangaDexID   string  `json:"mangadexId,omitempty"`
	AutoSync     bool    `json:"autoSync"`
	LastSyncedAt string  `json:"lastSyncedAt,omitempty"`
	LastChapter  float64 `json:"lastChapter,omitempty"`
	LastStatus   string  `json:"lastStatus,omitempty"`
}

type UserTrackersDocument struct {
	ID        string              `json:"_id,omitempty"`
	Type      string              `json:"_type,omitempty"`
	UpdatedAt string              `json:"updatedAt,omitempty"`
	Trackers  []SeriesTrackerItem `json:"trackers"`
}

type TrackerStatus struct {
	Provider        string  `json:"provider"`
	MatchedID       string  `json:"matchedId,omitempty"`
	Title           string  `json:"title,omitempty"`
	CoverImage      string  `json:"coverImage,omitempty"`
	Status          string  `json:"status,omitempty"` // reading, completed, on_hold, dropped, plan_to_read, rereading
	Score           float64 `json:"score,omitempty"`
	NumReadChapters float64 `json:"numReadChapters"`
	TotalChapters   int     `json:"totalChapters,omitempty"`
	URL             string  `json:"url,omitempty"`
	IsConnected     bool    `json:"isConnected"`
}

type SeriesAllTrackersStatus struct {
	SeriesID string                    `json:"seriesId"`
	AutoSync bool                      `json:"autoSync"`
	Trackers map[string]*TrackerStatus `json:"trackers"`
}

type TrackerSearchResult struct {
	ID            string `json:"id"`
	Provider      string `json:"provider"`
	Title         string `json:"title"`
	RomajiTitle   string `json:"romajiTitle,omitempty"`
	CoverImage    string `json:"coverImage,omitempty"`
	Format        string `json:"format,omitempty"`
	Status        string `json:"status,omitempty"`
	TotalChapters int    `json:"totalChapters,omitempty"`
	Year          int    `json:"year,omitempty"`
	URL           string `json:"url,omitempty"`
}

type BindTrackerRequest struct {
	LogtoID     string `json:"logtoId"`
	SeriesID    string `json:"seriesId"`
	SeriesTitle string `json:"seriesTitle,omitempty"`
	Provider    string `json:"provider"`
	MatchedID   string `json:"matchedId"`
	GlobalBind  bool   `json:"globalBind,omitempty"`
}

type UpdateTrackerStatusRequest struct {
	LogtoID  string  `json:"logtoId"`
	SeriesID string  `json:"seriesId"`
	Provider string  `json:"provider"`
	Status   string  `json:"status"`
	Chapter  float64 `json:"chapter"`
	Score    float64 `json:"score"`
}

type SyncProgressRequest struct {
	LogtoID       string  `json:"logtoId"`
	SeriesID      string  `json:"seriesId"`
	ChapterNumber float64 `json:"chapterNumber"`
	TotalChapters int     `json:"totalChapters,omitempty"`
}
