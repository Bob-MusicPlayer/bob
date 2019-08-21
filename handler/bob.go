package handler

import (
	"bob/util"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//{"event": "status", "payload": null}

type BobHandler struct {
	hub    *util.Hub
	worker *util.Worker
}

func NewBobHandler(hub *util.Hub, worker *util.Worker) *BobHandler {
	return &BobHandler{
		hub:    hub,
		worker: worker,
	}
}

func (bob *BobHandler) HandleConnect(w http.ResponseWriter, req *http.Request) {
	header := http.Header{}
	header.Set("Access-Control-Allow-Origin", "*")

	conn, err := upgrader.Upgrade(w, req, header)
	if err != nil {
		fmt.Println(err)
	}
	util.NewChannel(bob.hub, bob.worker, conn)
}
