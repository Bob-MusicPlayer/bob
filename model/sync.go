package model

type Sync struct {
	IsPlaying bool      `json:"isPlaying"`
	Playback  *Playback `json:"playback"`
}
