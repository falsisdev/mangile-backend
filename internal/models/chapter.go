package models

type Chapter struct {
	MyAnimeListId int    `json:"myAnimeListId"`
	Title         string `json:"title"`
	Type          string `json:"_type"`
	Chapter       any    `json:"chapter"`
}
