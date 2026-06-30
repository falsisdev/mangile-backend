package models

type AniListListResponse struct {
	Data AniListListData `json:"data"`
}

type AniListListData struct {
	Page AniListListPage `json:"Page"`
}

type AniListListPage struct {
	Media []AniListListMedia `json:"media"`
}

type AniListListMedia struct {
	ID                 int               `json:"id"`
	IDMal              int               `json:"idMal"`
	Title              AniListMangaTitle `json:"title"`
	Type               string            `json:"type"`
	Format             string            `json:"format"`
	Status             string            `json:"status"`
	MeanScore          int               `json:"meanScore"`
	Description		   string            `json:"description"`
	CoverImage         struct {
		Large string `json:"large"`
	} 									 `json:"coverImage"`
	BannerImage string 					 `json:"bannerImage"`
}

type MangaCard struct {
	AniListID          int    `json:"anilist_id"`
	MyAnimeListID      int    `json:"mal_id"`
	AniListTitle       string `json:"anilist_title"`
	TitleRomaji        string `json:"anilist_title_romaji"`
	TitleEnglish       string `json:"title_english"`
	TitleNative        string `json:"title_native"`
	Type               string `json:"anilist_type"`
	Format             string `json:"anilist_format"`
	Status             string `json:"anilist_status"`
	Score              int    `json:"anilist_score"`
	CoverImage         string `json:"anilist_cover_image"`
	BannerImage        string `json:"anilist_banner_image"`
	AniListDescription string `json:"anilist_description"`
}
