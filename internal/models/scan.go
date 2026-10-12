package models

type Scan struct {
	ID          string        `json:"_id"`
	Type        string        `json:"_type"`
	Logo        string        `json:"logo"`
	CoverImage  string        `json:"coverImage"`
	Members     []interface{} `json:"members"`
	Name        string        `json:"name"`
	Description []interface{} `json:"description"`
	Website     string        `json:"website"`
}
