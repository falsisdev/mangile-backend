package models

type SanityList struct {
	ID           string   `json:"_id"`
	Type         string   `json:"_type"`
	CreatedAt    string   `json:"_createdAt"`
	UpdatedAt    string   `json:"_updatedAt"`
	MalID        int      `json:"myAnimeListId"`
	SanityBanner string   `json:"bannerImage"`
	SanityCover  string   `json:"coverImage"`
	SanityTitle  string   `json:"title"`
	SanityTags   []string `json:"tags"`
}
