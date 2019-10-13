package player

import (
	"bob/model"
)

type Queue struct {
	Playbacks         []model.Playback
	PreviousPlaybacks []model.Playback
}

func NewQueue() *Queue {
	return &Queue{
		Playbacks: make([]model.Playback, 0),
	}
}

func (q *Queue) Size() int {
	return len(q.Playbacks)
}

func (q *Queue) SizePrevious() int {
	return len(q.PreviousPlaybacks)
}

func (q *Queue) NextAvailable() bool {
	return q.Size() > 0
}

func (q *Queue) PreviousAvailable() bool {
	return q.SizePrevious() > 0
}

func (q *Queue) AddPrevious(playback model.Playback) {
	q.PreviousPlaybacks = append([]model.Playback{playback}, q.PreviousPlaybacks...)
}

func (q *Queue) AddNext(playback model.Playback) {
	if q.Size() == 0 {
		q.PrependPlayback(playback)
	} else {
		q.insertIntoPlaybacks(1, playback)
	}
}

func (q *Queue) PrependPlayback(playback model.Playback) {
	q.Playbacks = append([]model.Playback{playback}, q.Playbacks...)
}

func (q *Queue) AppendPlayback(playback model.Playback) {
	q.Playbacks = append(q.Playbacks, playback)
}

func (q *Queue) RemoveFirst() {
	q.Playbacks = append(q.Playbacks[:0], q.Playbacks[1:]...)
}

func (q *Queue) RemoveFirstFromPrevious() {
	q.PreviousPlaybacks = append(q.PreviousPlaybacks[:0], q.PreviousPlaybacks[1:]...)
}

func (q *Queue) RemoveFirstAndAddToPrevious() {
	q.AddPrevious(q.Playbacks[0])
	q.Playbacks = append(q.Playbacks[:0], q.Playbacks[1:]...)
}

func (q *Queue) Clear() {
	q.Playbacks = make([]model.Playback, 0)
}

func (q *Queue) insertIntoPlaybacks(index int, item model.Playback) {
	q.Playbacks = append(q.Playbacks, model.Playback{})
	copy(q.Playbacks[index+1:], q.Playbacks[index:])
	q.Playbacks[index] = item
}
