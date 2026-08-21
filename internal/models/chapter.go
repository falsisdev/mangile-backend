package models

type ChapterInfo struct {
	Title         string `json:"title"`
	Key           string `json:"_key"`
	Type          string `json:"_type"`
	ChapterNumber int    `json:"chapterNumber"`
	Content       any    `json:"content"`
	Source        any    `json:"source"`
}

type Chapter struct {
	MyAnimeListId int         `json:"myAnimeListId"`
	Title         string      `json:"title"`
	Type          string      `json:"_type"`
	Chapter       ChapterInfo `json:"chapter"`
	ChapterKeys   []string    `json:"chapterKeys"`
}
