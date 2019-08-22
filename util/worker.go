package util

import (
	"bob/core"
	"bob/model"
	"bob/player"
	"fmt"
)

type Worker struct {
	hub    *Hub
	player *player.Player
	env    *core.Environment
}

func NewWorker(hub *Hub, player *player.Player, env *core.Environment) *Worker {
	return &Worker{
		hub:    hub,
		player: player,
		env:    env,
	}
}

func (w *Worker) HandleClientRequest(channel *Channel, packet *Packet) {
	switch packet.Event {
	case "setPlayback":
		w.handleSetPlayback(channel)
		break
	case "search":
		w.handleSearch(channel, packet.Payload.(string))
		break
	case "getPlayer":
		w.handleGetPlayer(channel)
	case "echo":
		channel.Send <- packet
		break
	default:
		channel.Send <- &Packet{
			Event:   "error",
			Payload: fmt.Sprintf("No event with name %s found", packet.Event),
		}
	}
}

func (w *Worker) handleSetPlayback(channel *Channel) {
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
}

func (w *Worker) handleSearch(channel *Channel, query string) {
	searchResponse := w.player.Search(query)

	channel.Send <- &Packet{
		Event:   "success",
		Payload: searchResponse,
	}
}

func (w *Worker) handleGetPlayer(channel *Channel) {
	channel.Send <- &Packet{
		Event:   "success",
		Payload: w.env.ConfigManager.Config.Player,
	}
}
