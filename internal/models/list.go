package models

type List struct {
	ID        string        `json:"_id"`
	Type      string        `json:"_type"` //lists
	CreatedAt string        `json:"createdAt"`
	Items     interface{}   `json:"items"`
	Title     string        `json:"title"`
	Likes     []interface{} `json:"likes"`
	User      interface{}   `json:"user"`
}
