package enet

/*
#cgo LDFLAGS: -lenet
#include <enet/enet.h>

// Helper to create ENet host
static ENetHost* create_host(const char* addr, int port, int peerCount) {
    ENetAddress address;
    enet_address_set_host(&address, addr);
    address.port = port;
    return enet_host_create(&address, peerCount, 2, 0, 0);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Host wraps an ENet host for UDP game communication.
type Host struct {
	cHost   *C.ENetHost
	peers   map[uint32]*Peer
	onConnect func(uint32)
	onReceive func(uint32, []byte)
	onDisconnect func(uint32)
}

// Init initializes the ENet library. Call once at startup.
func Init() error {
	if C.enet_initialize() != 0 {
		return fmt.Errorf("enet_initialize failed")
	}
	return nil
}

// Deinit cleans up the ENet library. Call at shutdown.
func Deinit() {
	C.enet_deinitialize()
}

// NewHost creates an ENet host bound to the given address and port.
func NewHost(addr string, port int, peerCount int) (*Host, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	cAddr := C.CString(addr)
	defer C.free(unsafe.Pointer(cAddr))

	cHost := C.create_host(cAddr, C.int(port), C.int(peerCount))
	if cHost == nil {
		return nil, fmt.Errorf("enet_host_create failed for %s:%d", addr, port)
	}

	return &Host{
		cHost: cHost,
		peers: make(map[uint32]*Peer),
	}, nil
}

// Service polls ENet for incoming events. Returns immediately.
// Call in a tight loop or with a 0 timeout for non-blocking.
func (h *Host) Service(timeout int) {
	var event C.ENetEvent
	for C.enet_host_service(h.cHost, &event, C.uint(timeout)) > 0 {
		peerID := uint32(event.peer.connectID)

		switch event._type {
		case C.ENET_EVENT_TYPE_CONNECT:
			h.peers[peerID] = &Peer{
				id:    peerID,
				cPeer: event.peer,
			}
			if h.onConnect != nil {
				h.onConnect(peerID)
			}

		case C.ENET_EVENT_TYPE_RECEIVE:
			data := C.GoBytes(unsafe.Pointer(event.packet.data), C.int(event.packet.dataLength))
			C.enet_packet_destroy(event.packet)
			if h.onReceive != nil {
				h.onReceive(peerID, data)
			}

		case C.ENET_EVENT_TYPE_DISCONNECT:
			delete(h.peers, peerID)
			if h.onDisconnect != nil {
				h.onDisconnect(peerID)
			}
		}
	}
}

// Send sends a packet to a peer (reliable).
func (h *Host) Send(peerID uint32, data []byte) error {
	p, ok := h.peers[peerID]
	if !ok {
		return fmt.Errorf("peer %d not found", peerID)
	}
	return p.Send(data)
}

// SendAll sends a packet to all connected peers.
func (h *Host) SendAll(data []byte) {
	for _, p := range h.peers {
		p.Send(data)
	}
}

// BroadcastExcept sends to all peers except the one specified.
func (h *Host) BroadcastExcept(exceptID uint32, data []byte) {
	for id, p := range h.peers {
		if id != exceptID {
			p.Send(data)
		}
	}
}

// Disconnect disconnects a peer gracefully.
func (h *Host) Disconnect(peerID uint32) {
	if p, ok := h.peers[peerID]; ok {
		p.Disconnect()
		delete(h.peers, peerID)
	}
}

// Close shuts down the host and frees all resources.
func (h *Host) Close() {
	for _, p := range h.peers {
		p.Disconnect()
	}
	h.peers = nil
	C.enet_host_destroy(h.cHost)
}

// SetOnConnect sets the connect callback.
func (h *Host) SetOnConnect(fn func(uint32)) { h.onConnect = fn }

// SetOnReceive sets the receive callback.
func (h *Host) SetOnReceive(fn func(uint32, []byte)) { h.onReceive = fn }

// SetOnDisconnect sets the disconnect callback.
func (h *Host) SetOnDisconnect(fn func(uint32)) { h.onDisconnect = fn }

// PeerCount returns the number of connected peers.
func (h *Host) PeerCount() int {
	return len(h.peers)
}
