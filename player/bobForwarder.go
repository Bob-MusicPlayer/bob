package player

import (
	"bob/core"
	"bob/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BobForwarder struct {
	env *core.Environment
}

func NewBoxForwarder(env *core.Environment) *BobForwarder {
	return &BobForwarder{
		env: env,
	}
}

func (bf *BobForwarder) buildUrl(player *model.Player, action string) string {
	return fmt.Sprintf("http://%s:%d/api/v1/%s", player.Hostname, player.Port, action)
}

func (bf *BobForwarder) ForwardSetPlayback(playback model.Playback) error {
	player := bf.env.ConfigManager.GetPlayerBySource(playback.Source)

	fmt.Println(bf.buildUrl(player, "playback"))

	data, err := json.Marshal(playback)
	if err != nil {
		return err
	}

	_, err = http.Post(bf.buildUrl(player, "playback"), "application/json", bytes.NewBuffer(data))

	return err
}
