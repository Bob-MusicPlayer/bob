package handler

import (
	"bob/model"
	"bob/player"
	"bob/utils"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"net/http"
	shared "shared-bob"
	"time"
)

type BobHandler struct {
	player      *player.Player
	eventBroker *sse.Server
}

func NewBobHandler(player *player.Player, eventBroker *sse.Server) *BobHandler {

	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			err := player.Sync()
			if err != nil {
				fmt.Println(fmt.Sprintf("Sync failed with error: %s", err.Error()))
			}

			sync := model.Sync{
				IsPlaying: player.IsPlaying,
				Playback:  player.CurrentPlayback,
			}

			_ = utils.SendEvent(eventBroker, "sync", sync)
		}
	}()

	return &BobHandler{
		player:      player,
		eventBroker: eventBroker,
	}
}

func (bh *BobHandler) HandlePlay(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)

	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	err := bh.player.Play()
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = utils.SendEvent(bh.eventBroker, "play", nil)
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(nil)
}

func (bh *BobHandler) HandlePause(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)

	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	err := bh.player.Pause()
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = utils.SendEvent(bh.eventBroker, "pause", nil)
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(nil)
}

func (bh *BobHandler) HandlePlayback(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)

	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	playback := model.Playback{
		ID:     "DKnIpsHe_YM",
		Source: "youtube",
	}

	err := bh.player.SetPlayback(playback)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = utils.SendEvent(bh.eventBroker, "play", nil)
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(nil)
}

func (bh *BobHandler) HandleSearch(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)

	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	var searchRequest model.SearchRequest

	err := responseHelper.DecodeBody(&searchRequest)
	if responseHelper.ReturnHasError(err) {
		return
	}

	searchResponse := bh.player.Search(&searchRequest)

	err = utils.SendEvent(bh.eventBroker, "play", nil)
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(searchResponse)
}
