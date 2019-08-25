package player

import (
	"bob/core"
	"bob/model"
	"fmt"
)

type Player struct {
	CurrentPlayback *model.Playback
	IsPlaying       bool
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

func (p *Player) Search(search *model.SearchRequest) *model.SearchResponse {
	return p.bobForwarder.ForwardSearch(search.Query)
}

func (p *Player) SetPlayback(playback model.Playback) error {
	p.Queue.Clear()
	p.Queue.PrependPlayback(playback)

	p.CurrentPlayback = &playback

	err := p.bobForwarder.ForwardSetPlayback(playback)

	return err
}

func (p *Player) Play() error {
	if p.CurrentPlayback == nil {
		return nil
	}
	err := p.bobForwarder.ForwardPlay(p.CurrentPlayback.Source)
	return err
}

func (p *Player) Pause() error {
	if p.CurrentPlayback == nil {
		return nil
	}
	err := p.bobForwarder.ForwardPause(p.CurrentPlayback.Source)
	return err
}

func (p *Player) Sync() error {
	if p.CurrentPlayback == nil {
		return nil
	}
	playback, err := p.bobForwarder.ForwardGetPlaybackInfo(p.CurrentPlayback.Source)
	if err != nil {
		return err
	}

	p.IsPlaying = playback.IsPlaying

	p.CurrentPlayback = playback
	return nil
}

func (p *Player) SeekTo(source string, seconds int) error {
	if p.CurrentPlayback == nil {
		return nil
	}

	fmt.Println(seconds)

	err := p.bobForwarder.ForwardSeek(source, seconds)
	if err != nil {
		return err
	}

	return nil
}
