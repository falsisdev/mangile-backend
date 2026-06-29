package models

type User struct {
	ID          string        `json:"_id"`
	LogtoID     string        `json:"logtoId"`
	Type        string        `json:"_type"` //auth
	CreatedAt   string        `json:"_createdAt"`
	Avatar      string        `json:"avatar"`
	BannerImage string        `json:"banner"`
	FavChapters []interface{} `json:"favoriteChapters"`
	Biography   string        `json:"bio"`
	FavScans    []interface{} `json:"favoriteScans"`
	FavTitle    string        `json:"favoriteTitle"`
	FavTitles   []string      `json:"favoriteTitles"`
	Gender      string        `json:"gender"`
	Lists       []interface{} `json:"lists"`
	Name        string        `json:"name"`
	Username    string        `json:"username"`
}

//Bu sürümde bookcase nesnesi tamamen kaldırıldı
//Onun yerine Anilist/MAL Sync eklenecek ilerleyen dönemlerde
