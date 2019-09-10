package model

type QueueEntry struct {
	ID         int `sql:",pk,type:uuid default uuid_generate_v4()"`
	Index      int
	PlaybackID int `sql:"type:uuid,notnull"`
}
