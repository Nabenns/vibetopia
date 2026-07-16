package enet

import (
	"strings"
)

// Packet represents a parsed Growtopia game packet.
// Format: pipe-delimited key|value\n pairs.
type Packet map[string]string

// ParsePacket parses a raw byte slice into a Packet.
func ParsePacket(data []byte) Packet {
	p := Packet{}
	text := strings.TrimRight(string(data), "\x00")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		idx := strings.IndexByte(line, '|')
		if idx < 0 {
			continue
		}
		key := line[:idx]
		val := line[idx+1:]
		p[key] = val
	}
	return p
}

// Get returns a value from the packet, or empty string if not present.
func (p Packet) Get(key string) string {
	return p[key]
}

// Serialize serializes a Packet back to bytes for sending over ENet.
func (p Packet) Serialize() []byte {
	var b strings.Builder
	for k, v := range p {
		b.WriteString(k)
		b.WriteByte('|')
		b.WriteString(v)
		b.WriteByte('\n')
	}
	return []byte(b.String())
}

// VarList builds a VariantList-style packet from a map of strings.
// VariantList format: key|value\nkey|value\n (same as Packet, but with onConsoleMessage prefix)
func VarList(items map[string]string) []byte {
	p := Packet(items)
	return p.Serialize()
}

// Quick helpers for common packet building

// OnConsoleMessage builds an onConsoleMessage packet.
func OnConsoleMessage(msg string) Packet {
	return Packet{
		"action":            "onConsoleMessage",
		"onConsoleMessage":  msg,
	}
}

// OnDialogRequest builds an onDialogRequest packet.
func OnDialogRequest(dialog string) Packet {
	return Packet{
		"action":          "onDialogRequest",
		"onDialogRequest": dialog,
	}
}

// LogOn sends a log entry to the client console.
func (h *Host) LogOn(peerID uint32, msg string) error {
	return h.Send(peerID, OnConsoleMessage(msg).Serialize())
}

// Dialog shows a dialog to the player.
func (h *Host) Dialog(peerID uint32, name string) error {
	return h.Send(peerID, OnDialogRequest(name).Serialize())
}
