package handler

import (
	"bob/model"
	"bob/player"
	"bob/utils"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"net/http"
	shared "shared-bob"
	"strconv"
	"time"
)

type BobHandler struct {
	player      *player.Player
	syncEnabled bool
	eventBroker *sse.Server
}

func NewBobHandler(player *player.Player, eventBroker *sse.Server) *BobHandler {
	handler := &BobHandler{
		player:      player,
		eventBroker: eventBroker,
		syncEnabled: false,
	}

	handler.SendSyncs()

	return handler
}

func (bh *BobHandler) SendSyncs() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			if bh.syncEnabled {
				bh.sync()
			}
		}
	}()
}

func (bh *BobHandler) sync() {
	err := bh.player.Sync()
	if err != nil {
		fmt.Println(fmt.Sprintf("Sync failed with error: %s", err.Error()))
	}

	sync := model.Sync{
		PlayerState:       bh.player.GetState(),
		Playback:          bh.player.CurrentPlayback,
		NextAvailable:     bh.player.Queue.NextAvailable(),
		PreviousAvailable: bh.player.Queue.PreviousAvailable(),
	}

	_ = utils.SendEvent(bh.eventBroker, "sync", sync)
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

func (bh *BobHandler) HandleNext(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)

	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	err := utils.SendEvent(bh.eventBroker, "loading", true)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = bh.player.Next()
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = utils.SendEvent(bh.eventBroker, "loading", false)
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(nil)
}

func (bh *BobHandler) HandlePrevious(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)

	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	err := utils.SendEvent(bh.eventBroker, "loading", true)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = bh.player.Previous()
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = utils.SendEvent(bh.eventBroker, "loading", false)
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

	bh.syncEnabled = false
	bh.player.CurrentPlayback = nil

	err := utils.SendEvent(bh.eventBroker, "loading", true)
	if responseHelper.ReturnHasError(err) {
		return
	}
	bh.sync()

	var playback model.Playback

	err = responseHelper.DecodeBody(&playback)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = bh.player.SetPlayback(playback)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = utils.SendEvent(bh.eventBroker, "loading", false)
	if responseHelper.ReturnHasError(err) {
		return
	}
	bh.sync()

	bh.syncEnabled = true

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

	fmt.Println("test")

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

func (bh *BobHandler) HandlePlaybackSeek(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)

	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	seconds := req.URL.Query().Get("seconds")

	sec, err := strconv.Atoi(seconds)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = bh.player.SeekTo(bh.player.CurrentPlayback.Source, sec)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = utils.SendEvent(bh.eventBroker, "seek", sec)
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(nil)
}

func (bh *BobHandler) HandleSync(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)
	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	bh.sync()

	responseHelper.ReturnOk(nil)
}

func (bh *BobHandler) HandleQueueNext(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)

	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	var playback model.Playback

	err := responseHelper.DecodeBody(&playback)
	if responseHelper.ReturnHasError(err) {
		return
	}

	bh.player.Queue.AddNext(playback)

	bh.sync()

	responseHelper.ReturnOk(nil)
}
