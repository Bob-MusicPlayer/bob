package model

type Sync struct {
	PlayerState       string    `json:"playerState"`
	NextAvailable     bool      `json:"nextAvailable"`
	PreviousAvailable bool      `json:"previousAvailable"`
	Playback          *Playback `json:"playback"`
}
