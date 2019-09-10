package repository

import (
	"bob/model"
	"github.com/go-pg/pg/v9"
)

type PlayerRepository struct {
	database *pg.DB
}

func NewPlayerRepository(database *pg.DB) *PlayerRepository {
	return &PlayerRepository{
		database: database,
	}
}

func (br *PlayerRepository) PrependQueue()

func (br *PlayerRepository) AddPlayback(playback model.Playback) error {
	return br.database.Insert(playback)
}

func (br *PlayerRepository) GetPlayback(id string) (*model.Playback, error) {
	playback := new(model.Playback)
	err := br.database.Model(playback).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, err
	}

	return playback, nil
}
