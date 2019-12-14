package player

import (
	"errors"
	"github.com/Bob-MusicPlayer/bob/core"
	"github.com/Bob-MusicPlayer/bob/model"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	PLAYER_STATE_PLAYING     = "PLAYING"
	PLAYER_STATE_LOADING     = "LOADING"
	PLAYER_STATE_PAUSED      = "PAUSED"
	PLAYER_STATE_NO_PLAYBACK = "NOPLAYBACK"
)

var (
	NO_NEXT_PLAYBACK = errors.New("no next playback in queue available")
)

type Player struct {
	CurrentPlayback *model.Playback
	IsPlaying       bool
	Queue           *Queue
	loading         bool
	env             *core.Environment
	bobForwarder    *BobForwarder
	eventBroker     *sse.Server
}

func NewPlayer(queue *Queue, env *core.Environment, bobForwarder *BobForwarder, eventBroker *sse.Server) *Player {
	player := &Player{
		Queue:        queue,
		env:          env,
		bobForwarder: bobForwarder,
		eventBroker:  eventBroker,
	}

	player.PlayerLoop()

	return player
}

func (p *Player) PlayerLoop() {
	go func() {
		oldPosition := 0.0
		oldPosEqPosTimes := 0

		for {
			logrus.WithFields(logrus.Fields{
				"current":  p.CurrentPlayback,
				"position": p.IsPlaying,
			})
			if p.CurrentPlayback != nil {
				oldPosition = p.CurrentPlayback.Position
			}
			p.Sync()
			if p.CurrentPlayback != nil && p.CurrentPlayback.Position >= 0 && p.CurrentPlayback.Position == oldPosition {
				oldPosEqPosTimes++
			} else {
				oldPosEqPosTimes = 0
			}
			if p.IsPlaying && p.CurrentPlayback != nil && p.Queue.Size() == 0 && p.CurrentPlayback.Duration > 0 && oldPosEqPosTimes > 4 {
				logrus.Info("Current Playback finished. Stop player.")
				p.CurrentPlayback = nil
				oldPosEqPosTimes = 0
			} else if p.IsPlaying && p.CurrentPlayback != nil && p.Queue.Size() > 0 && p.CurrentPlayback.Duration > 0 && oldPosEqPosTimes > 4 {
				logrus.Info("Current Playback finished. Set next playback.")
				p.Next()
				oldPosEqPosTimes = 0
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()
}

func (p *Player) GetState() string {
	if p.loading {
		return PLAYER_STATE_LOADING
	}
	if p.CurrentPlayback == nil {
		return PLAYER_STATE_NO_PLAYBACK
	}
	if p.IsPlaying {
		return PLAYER_STATE_PLAYING
	} else {
		return PLAYER_STATE_PAUSED
	}
}

func (p *Player) Search(search *model.SearchRequest) *model.SearchResponse {
	logrus.WithField("search", search).Info("Search for Playbacks")
	return p.bobForwarder.ForwardSearch(search.Query)
}

func (p *Player) SetPlayback(playback model.Playback) error {
	p.loading = true

	logrus.WithField("playback", playback).Info("Set playback")

	err := p.bobForwarder.ForwardSetPlayback(playback)

	p.CurrentPlayback = &playback

	p.loading = false

	return err
}

func (p *Player) Play() error {
	if p.CurrentPlayback == nil {
		return nil
	}

	logrus.Info("Play")

	err := p.bobForwarder.ForwardPlay(p.CurrentPlayback.Source)
	return err
}

func (p *Player) Pause() error {
	if p.CurrentPlayback == nil {
		return nil
	}

	logrus.Info("Pause")

	err := p.bobForwarder.ForwardPause(p.CurrentPlayback.Source)
	return err
}

func (p *Player) Next() error {
	logrus.Info("Next Playback")

	if p.Queue.Size() == 0 {
		return NO_NEXT_PLAYBACK
	}

	p.Queue.AddPrevious(*p.CurrentPlayback)

	err := p.SetPlayback(p.Queue.Playbacks[0])
	if err != nil {
		return err
	}

	p.Queue.RemoveFirst()

	return nil
}

func (p *Player) Previous() error {
	logrus.Info("Previous Playback")

	if p.Queue.SizePrevious() == 0 {
		return NO_NEXT_PLAYBACK
	}

	p.Queue.PrependPlayback(*p.CurrentPlayback)

	err := p.SetPlayback(p.Queue.PreviousPlaybacks[0])
	if err != nil {
		return err
	}

	p.Queue.RemoveFirstFromPrevious()

	return nil
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
	logrus.WithField("seconds", seconds).Info("Seek")

	if p.CurrentPlayback == nil {
		return nil
	}

	err := p.bobForwarder.ForwardSeek(source, seconds)
	if err != nil {
		return err
	}

	return nil
}
