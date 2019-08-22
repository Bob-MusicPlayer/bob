package model

type SearchResponse map[string]PlayerSearchResponse

type PlayerSearchResponse struct {
	Amount    int        `json:"amount"`
	Error     string     `json:"error"`
	Playbacks []Playback `json:"playbacks"`
}
