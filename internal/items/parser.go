package items

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const xorToken = "PBG892FXX982ABC*"

type Item struct {
	ID   uint32
	Name string
}

type Store struct {
	Version uint16
	Items   map[uint32]*Item
	IDs     []uint32
}

func Parse(r io.Reader) (*Store, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)

	var version uint16
	var count uint32
	binary.Read(buf, binary.LittleEndian, &version)
	binary.Read(buf, binary.LittleEndian, &count)

	store := &Store{
		Version: version,
		Items:   make(map[uint32]*Item, count),
		IDs:     make([]uint32, 0, count),
	}

	for i := uint32(0); i < count; i++ {
		item, err := parseItemCompact(buf, version)
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		store.Items[item.ID] = item
		store.IDs = append(store.IDs, item.ID)
	}

	return store, nil
}

func parseItemCompact(r *bytes.Reader, version uint16) (*Item, error) {
	// ID: read 4 bytes (Gurotopia reads id as uint32 then skips 2)
	var id uint32
	binary.Read(r, binary.LittleEndian, &id)
	skip(r, 2)

	skip(r, 4) // property, cat, type, pad

	// Name: short length-prefixed, XOR encoded
	var nameLen uint16
	binary.Read(r, binary.LittleEndian, &nameLen)
	nameBytes := make([]byte, nameLen)
	io.ReadFull(r, nameBytes)
	for j := range nameBytes {
		nameBytes[j] ^= xorToken[(j+int(id))%len(xorToken)]
	}
	name := cstr(nameBytes)

	// filename (length-prefixed short string)
	readSkipStr(r) // filename

	// skip: int + byte + ingredient + 4 + collision + hits + hit_reset
	skip(r, 4+1+1+4+4+4+4)

	// clothing type or skip
	skip(r, 1)

	// rarity
	skip(r, 2)

	// skip byte
	skip(r, 1)

	// audio filename (length-prefixed)
	readSkipStr(r)
	skip(r, 4) // audio result

	skip(r, 4) // unknown

	// 4x short-length-prefixed strings
	readSkipStr(r)
	readSkipStr(r)
	readSkipStr(r)
	readSkipStr(r)

	skip(r, 16) // array<u_char, 16>

	skip(r, 4) // tick
	skip(r, 2+2) // 2x short

	// 2x short-length-prefixed strings
	readSkipStr(r)
	readSkipStr(r)

	skip(r, 80) // array<u_char, 80>

	// Version-dependent fields
	if version >= 0x0b { readSkipStr(r) }
	if version >= 0x0c { skip(r, 4+9) }
	if version >= 0x0d { skip(r, 4) }
	if version >= 0x0e { skip(r, 4) }
	if version >= 0x0f { skip(r, 25); readSkipStr(r) }
	if version >= 0x10 { readSkipStr(r) }
	if version >= 0x11 { skip(r, 4) }
	if version >= 0x12 { skip(r, 4) }
	if version >= 0x13 { skip(r, 9) }
	if version >= 0x15 { skip(r, 2) }
	if version >= 0x16 { readSkipStr(r) }
	if version >= 0x17 { skip(r, 4+4) }
	if version >= 0x18 { skip(r, 1) }
	if version >= 0x19 { readSkipStr(r); skip(r, 4) }
	if version == 0x1a { skip(r, 1) }

	return &Item{ID: id, Name: name}, nil
}

func readSkipStr(r *bytes.Reader) int {
	var length uint16
	binary.Read(r, binary.LittleEndian, &length)
	skip(r, int(length))
	return int(length)
}

func skip(r *bytes.Reader, n int) {
	r.Seek(int64(n), io.SeekCurrent)
}

func cstr(b []byte) string {
	end := bytes.IndexByte(b, 0)
	if end < 0 {
		return string(b)
	}
	return string(b[:end])
}

func (s *Store) Get(id uint32) *Item { return s.Items[id] }
func (s *Store) GetName(id uint32) string {
	if item, ok := s.Items[id]; ok {
		return item.Name
	}
	return fmt.Sprintf("Unknown(%d)", id)
}
