package model

type Playback struct {
	DBID          int     `sql:"id,pk,type:uuid default uuid_generate_v4()" json:"-"`
	ID            string  `sql:"source_id,nopk" json:"id"`
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Position      float64 `sql:"-" json:"position"`
	Duration      float64 `sql:"-" json:"duration"`
	CachePosition float64 `sql:"-" json:"cachePosition"`
	Source        string  `json:"source"`
	ThumbnailUrl  string  `json:"thumbnailUrl"`
	IsPlaying     bool    `sql:"-" json:"isPlaying"`
}
