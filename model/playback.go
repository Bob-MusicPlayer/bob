package model

type Playback struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Position      float64 `json:"position"`
	Duration      float64 `json:"duration"`
	CachePosition float64 `json:"cachePosition"`
	Source        string  `json:"source"`
	ThumbnailUrl  string  `json:"thumbnailUrl"`
}
