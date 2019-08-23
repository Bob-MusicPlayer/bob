package model

type SearchRequest struct {
	Source string `json:"source"`
	Query  string `json:"query"`
}
