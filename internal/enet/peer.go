package enet

/*
#cgo LDFLAGS: -lenet
#include <enet/enet.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Peer wraps an ENet peer connection.
type Peer struct {
	id    uint32
	cPeer *C.ENetPeer
}

// ID returns the peer's connection ID.
func (p *Peer) ID() uint32 { return p.id }

// Send sends a reliable packet to this peer.
func (p *Peer) Send(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	flags := C.uint(C.ENET_PACKET_FLAG_RELIABLE)
	packet := C.enet_packet_create(unsafe.Pointer(&data[0]), C.size_t(len(data)), flags)
	if packet == nil {
		return fmt.Errorf("enet_packet_create failed")
	}
	if C.enet_peer_send(p.cPeer, 0, packet) < 0 {
		return fmt.Errorf("enet_peer_send failed for peer %d", p.id)
	}
	return nil
}

// Disconnect gracefully disconnects this peer.
func (p *Peer) Disconnect() {
	C.enet_peer_disconnect(p.cPeer, 0)
}
