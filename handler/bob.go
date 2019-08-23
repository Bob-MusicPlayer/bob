package handler

import (
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
}

func NewBobHandler() *BobHandler {
	return &BobHandler{}
}
