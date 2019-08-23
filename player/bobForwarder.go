package player

import (
	"bob/core"
	"bob/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (bf *BobForwarder) getPlayer() []*model.Player {
	return bf.env.ConfigManager.Config.Player
}

func (bf *BobForwarder) ForwardSetPlayback(playback model.Playback) error {
	player := bf.env.ConfigManager.GetPlayerBySource(playback.Source)

	data, err := json.Marshal(playback)
	if err != nil {
		return err
	}

	_, err = http.Post(bf.buildUrl(player, "playback"), "application/json", bytes.NewBuffer(data))

	return err
}

func (bf *BobForwarder) ForwardSearch(query string) *model.SearchResponse {
	players := bf.getPlayer()

	response := model.SearchResponse{}

	for _, player := range players {
		resp, err := http.Get(bf.buildUrl(player, fmt.Sprintf("search?q=%s", query)))
		if err != nil {
			response[player.Source] = model.PlayerSearchResponse{
				Amount:    0,
				Error:     err.Error(),
				Playbacks: make([]model.Playback, 0),
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {
			response[player.Source] = model.PlayerSearchResponse{
				Amount:    0,
				Error:     fmt.Sprintf("Request failed with Code %d", resp.StatusCode),
				Playbacks: make([]model.Playback, 0),
			}
			continue
		}
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			response[player.Source] = model.PlayerSearchResponse{
				Amount:    0,
				Error:     err.Error(),
				Playbacks: make([]model.Playback, 0),
			}
			continue
		}
		var playbacks []model.Playback
		err = json.Unmarshal(raw, &playbacks)
		if err != nil {
			response[player.Source] = model.PlayerSearchResponse{
				Amount:    0,
				Error:     err.Error(),
				Playbacks: make([]model.Playback, 0),
			}
			continue
		}
		response[player.Source] = model.PlayerSearchResponse{
			Amount:    len(playbacks),
			Error:     "",
			Playbacks: playbacks,
		}
	}
	return &response
}

func (bf *BobForwarder) ForwardPlay(source string) error {
	player := bf.env.ConfigManager.GetPlayerBySource(source)

	_, err := http.Post(bf.buildUrl(player, "play"), "application/json", nil)

	return err
}

func (bf *BobForwarder) ForwardPause(source string) error {
	player := bf.env.ConfigManager.GetPlayerBySource(source)

	_, err := http.Post(bf.buildUrl(player, "pause"), "application/json", nil)

	return err
}

func (bf *BobForwarder) ForwardGetPlaybackInfo(source string) (*model.Playback, error) {
	player := bf.env.ConfigManager.GetPlayerBySource(source)

	resp, err := http.Get(bf.buildUrl(player, "playback/info"))
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var playback model.Playback

	err = json.Unmarshal(raw, &playback)
	if err != nil {
		return nil, err
	}

	return &playback, err
}
