package util

type Worker struct {
	hub *Hub
}

func NewWorker(hub *Hub) *Worker {
	return &Worker{
		hub: hub,
	}
}

func (w *Worker) HandleClientRequest(channel *Channel, packet *Packet) {
	switch packet.Event {
	case "broadcast":
		w.hub.Broadcast(packet)
		break
	case "echo":
		channel.Send <- packet
		break
	}
}
