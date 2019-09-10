package utils

import (
	model2 "bob/model"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func InitializePlayerTables(db *pg.DB) error {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		return err
	}

	err = createPlaybackTable(db)
	if err != nil {
		return err
	}

	err = createQueueTable(db)
	if err != nil {
		return err
	}

	return nil
}

func createQueueTable(db *pg.DB) error {
	for _, model := range []interface{}{(*model2.QueueEntry)(nil), (*model2.QueueEntry)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func createPlaybackTable(db *pg.DB) error {
	for _, model := range []interface{}{(*model2.Playback)(nil), (*model2.Playback)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
