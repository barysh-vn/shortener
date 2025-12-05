package model

type Link struct {
	URL       string `json:"url"`
	Alias     string `json:"alias"`
	UserID    string `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
}
