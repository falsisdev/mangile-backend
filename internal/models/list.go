package models

type ListItemReference struct {
	ID            string      `json:"_id"`
	Type          string      `json:"_type"`
	Title         string      `json:"title"`
	Slug          string      `json:"slug"`
	CoverImage    *ImageAsset `json:"coverImage,omitempty"`
	BannerImage   *ImageAsset `json:"bannerImage,omitempty"`
	MyAnimeListID int         `json:"myAnimeListId,omitempty"`
	UploadStatus  string      `json:"uploadStatus,omitempty"`
	Format        string      `json:"format,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
}

type ListItem struct {
	Key          string             `json:"_key"`
	Reference    *ListItemReference `json:"reference,omitempty"`
	ExternalID   *int               `json:"externalId,omitempty"`
	ExternalType string             `json:"externalType,omitempty"`
	Index        int                `json:"index"`
	Note         string             `json:"note,omitempty"`
	AddedAt      string             `json:"addedAt,omitempty"`
}

type ListUser struct {
	ID       string `json:"id"`
	LogtoID  string `json:"logtoId"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

type List struct {
	ID          string      `json:"_id"`
	Type        string      `json:"_type"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	IsPrivate   bool        `json:"isPrivate"`
	CoverImage  *ImageAsset `json:"coverImage,omitempty"`
	CreatedAt   string      `json:"createdAt,omitempty"`
	UpdatedAt   string      `json:"updatedAt,omitempty"`
	User        *ListUser   `json:"user,omitempty"`
	Likes       []ListUser  `json:"likes,omitempty"`
	Items       []ListItem  `json:"items"`
	LikesCount  int         `json:"likesCount"`
	ItemsCount  int         `json:"itemsCount"`
}

type CreateListRequest struct {
	LogtoID     string `json:"logtoId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"isPrivate"`
}

type UpdateListRequest struct {
	LogtoID     string  `json:"logtoId"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	IsPrivate   *bool   `json:"isPrivate,omitempty"`
}

type AddListItemRequest struct {
	LogtoID      string  `json:"logtoId"`
	ReferenceID  *string `json:"referenceId,omitempty"`
	ExternalID   *int    `json:"externalId,omitempty"`
	ExternalType *string `json:"externalType,omitempty"`
	Index        *int    `json:"index,omitempty"`
	Note         string  `json:"note,omitempty"`
}

type UpdateListItemRequest struct {
	LogtoID string  `json:"logtoId"`
	Index   *int    `json:"index,omitempty"`
	Note    *string `json:"note,omitempty"`
}

type UpdateListItemsRequest struct {
	LogtoID string     `json:"logtoId"`
	Items   []ListItem `json:"items"`
}

type ToggleLikeListRequest struct {
	LogtoID string `json:"logtoId"`
}
