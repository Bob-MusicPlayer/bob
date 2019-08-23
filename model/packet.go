package model

type Packet struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}
