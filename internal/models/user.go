package models

type User struct {
	ID           string        `json:"_id"`
	LogtoID      string        `json:"logtoId"`
	Type         string        `json:"_type"` //auth
	CreatedAt    string        `json:"_createdAt"`
	Avatar       string        `json:"avatar"`
	BannerImage  string        `json:"banner"`
	Biography    string        `json:"bio"`
	Roles        []string      `json:"roles"`
	Lists        []interface{} `json:"lists"`
	Integrations []interface{} `json:"integrations"`
	Name         string        `json:"name"`
	Username     string        `json:"username"`
}
