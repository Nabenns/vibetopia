package protocol

import (
	"log"

	"vibetopia/internal/enet"
)

// Registry maps action names (like "tankIDName", "drop") to handler functions.
type Registry map[string]func(*Context)

// Context holds all state needed by a protocol handler.
type Context struct {
	PeerID uint32
	Packet enet.Packet
	Host   *enet.Host
	Player *PlayerState
}

// PlayerState holds runtime state for a connected player.
type PlayerState struct {
	GrowID  string
	InWorld bool
	World   string
	X, Y    int
}

// NewRegistry creates a new protocol registry.
func NewRegistry() Registry {
	return make(Registry)
}

// Handle dispatches a raw packet to the appropriate handler.
func (r Registry) Handle(peerID uint32, host *enet.Host, data []byte) {
	pkt := enet.ParsePacket(data)
	action := pkt.Get("action")
	log.Printf("[protocol] peer=%d action=%s", peerID, action)

	if handler, ok := r[action]; ok {
		ctx := &Context{
			PeerID: peerID,
			Packet: pkt,
			Host:   host,
			Player: &PlayerState{},
		}
		handler(ctx)
	} else {
		host.LogOn(peerID, "Unknown action: `"+action+"`")
	}
}
