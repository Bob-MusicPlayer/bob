package player

import (
	"bob/core"
	"bob/model"
)

type Player struct {
	CurrentPlayback *model.Playback
	IsPaused        bool
	Queue           *Queue
	env             *core.Environment
	bobForwarder    *BobForwarder
}

func NewPlayer(queue *Queue, env *core.Environment, bobForwarder *BobForwarder) *Player {
	return &Player{
		Queue:        queue,
		env:          env,
		bobForwarder: bobForwarder,
	}
}

func (p *Player) Search(query string) *model.SearchResponse {
	return p.bobForwarder.ForwardSearch(query)
}

func (p *Player) SetPlayback(playback model.Playback) error {
	p.Queue.Clear()
	p.Queue.PrependPlayback(playback)

	err := p.bobForwarder.ForwardSetPlayback(playback)

	return err
}
