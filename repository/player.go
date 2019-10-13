package repository

import (
	"bob/model"
)

type PlayerRepository struct {
}

func NewPlayerRepository() *PlayerRepository {
	return &PlayerRepository{}
}

func (br *PlayerRepository) PrependQueue() {

}

func (br *PlayerRepository) AddPlayback(playback model.Playback) error {

}

func (br *PlayerRepository) GetPlayback(id string) (*model.Playback, error) {
	playback := new(model.Playback)

	return playback, nil
}
