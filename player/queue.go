package player

import "bob/model"

type Queue struct {
	Playbacks []model.Playback
}

func NewQueue() *Queue {
	return &Queue{
		Playbacks: make([]model.Playback, 0),
	}
}

func (q *Queue) PrependPlayback(playback model.Playback) {
	q.Playbacks = append([]model.Playback{playback}, q.Playbacks...)
}

func (q *Queue) AppendPlayback(playback model.Playback) {
	q.Playbacks = append([]model.Playback{playback}, q.Playbacks...)
}

func (q *Queue) Clear() {
	q.Playbacks = make([]model.Playback, 0)
}
