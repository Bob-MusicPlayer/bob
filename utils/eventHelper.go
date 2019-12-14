package utils

import (
	"encoding/json"
	"github.com/Bob-MusicPlayer/bob/model"
	"github.com/alexandrevicenzi/go-sse"
)

func SendEvent(eventBroker *sse.Server, event string, payload interface{}) error {

	packet := model.Packet{
		Event:   event,
		Payload: payload,
	}

	data, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	eventBroker.SendMessage("/api/v1/events", sse.SimpleMessage(string(data)))

	return nil
}
