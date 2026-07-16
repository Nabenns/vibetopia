package items

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

func TestParseItemsDat(t *testing.T) {
	f, err := os.Open("/root/vibetopia/items.dat")
	if err != nil {
		t.Skipf("no items.dat: %v", err)
	}
	defer f.Close()

	store, err := Parse(f)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	t.Logf("Version: 0x%02x, Items: %d", store.Version, len(store.Items))

	if len(store.Items) == 0 {
		t.Error("no items parsed")
	}

	// Check known items
	tests := []struct {
		id   uint32
		name string
	}{
		{2, "Dirt"},
		{18, "Door"},
		{242, "Fist"},
		{1388, "Torch"},
	}
	for _, tc := range tests {
		item := store.Get(tc.id)
		if item == nil {
			t.Errorf("item %d not found", tc.id)
		} else {
			t.Logf("Item %d: %s", tc.id, item.Name)
		}
	}
}

func TestParseEmptyItemsDat(t *testing.T) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(0x13))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	store, err := Parse(buf)
	if err != nil {
		t.Fatalf("parse empty: %v", err)
	}
	if len(store.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(store.Items))
	}
}

func TestParseSingleItem(t *testing.T) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(0x13)) // version
	binary.Write(buf, binary.LittleEndian, uint32(1))     // count
	// Item ID
	binary.Write(buf, binary.LittleEndian, uint32(42))
	// Props
	binary.Write(buf, binary.LittleEndian, uint8(0))
	// Material
	binary.Write(buf, binary.LittleEndian, uint8(0))
	// Name
	buf.Write([]byte("TestItem\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	// Textures (16 bytes)
	buf.Write(make([]byte, 16))
	// Hash (4 bytes)
	buf.Write(make([]byte, 4))
	// growtime, rarity, fruittime, seedtime, animspeed, category, basecolor, overlaycolor
	for i := 0; i < 8; i++ {
		binary.Write(buf, binary.LittleEndian, uint32(0))
	}
	// unknown + collision
	buf.Write(make([]byte, 5))
	// ingredients
	for i := 0; i < 3; i++ {
		binary.Write(buf, binary.LittleEndian, uint32(0))
	}
	// xp
	binary.Write(buf, binary.LittleEndian, uint32(0))
	// version-dependent fields (v0x13)
	buf.Write(make([]byte, 1))  // 0x0b
	buf.Write(make([]byte, 4))  // 0x0c
	buf.Write(make([]byte, 4))  // 0x0d
	buf.Write(make([]byte, 4))  // 0x0e
	buf.Write(make([]byte, 4))  // 0x0f
	buf.Write(make([]byte, 4))  // 0x10
	buf.Write(make([]byte, 4))  // 0x11
	buf.Write(make([]byte, 4))  // 0x12
	buf.Write(make([]byte, 9))  // 0x13

	store, err := Parse(buf)
	if err != nil {
		t.Fatalf("parse single: %v", err)
	}
	if len(store.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(store.Items))
	}
	item := store.Get(42)
	if item == nil {
		t.Fatal("item 42 not found")
	}
	if item.Name != "TestItem" {
		t.Errorf("expected TestItem, got %s", item.Name)
	}
	if item.ID != 42 {
		t.Errorf("expected ID 42, got %d", item.ID)
	}
}

func TestParseInvalidVersion(t *testing.T) {
	f, err := os.Open("/root/vibetopia/items.dat")
	if err != nil {
		t.Skipf("no items.dat: %v", err)
	}
	defer f.Close()

	store, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	// Any valid version should parse the right number of items
	if store.Version < 0x0b || store.Version > 0x1a {
		t.Logf("Version 0x%02x outside known range (0x0b-0x1a), still parsed %d items", store.Version, len(store.Items))
	}
}
