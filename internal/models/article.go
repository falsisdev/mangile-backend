package models

type Article struct {
	ID          string      `json:"_id"`
	Type        string      `json:"_type"` //articles
	CreatedAt   string      `json:"_createdAt"`
	UpdatedAt   string      `json:"_updatedAt"`
	Article     interface{} `json:"article"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Slug        interface{} `json:"slug"`
}
