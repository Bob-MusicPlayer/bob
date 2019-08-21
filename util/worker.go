package util

import (
	"bob/model"
	"bob/player"
)

type Worker struct {
	hub    *Hub
	player *player.Player
}

func NewWorker(hub *Hub, player *player.Player) *Worker {
	return &Worker{
		hub:    hub,
		player: player,
	}
}

func (w *Worker) HandleClientRequest(channel *Channel, packet *Packet) {
	switch packet.Event {
	case "setPlayback":
		err := w.player.SetPlayback(model.Playback{
			ID:     "sIB7aDXKyyQ",
			Source: "youtube",
		})
		if err != nil {
			channel.Send <- &Packet{
				Event:   "error",
				Payload: err.Error(),
			}
			return
		}
		channel.Send <- &Packet{
			Event:   "success",
			Payload: nil,
		}
		break
	case "echo":
		channel.Send <- packet
		break
	}
}
