package uuid

import (
	"encoding/hex"
	"github.com/boostgo/lite/types/to"
	"github.com/google/uuid"
)

type UUID interface {
	String() string
}

type uid struct {
	id uuid.UUID
}

func New() UUID {
	return &uid{
		id: uuid.New(),
	}
}

func FromID(id uuid.UUID) UUID {
	return &uid{
		id: id,
	}
}

func Parse(input string) (UUID, error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}

	return FromID(id), nil
}

func ParseBytes(input []byte) (UUID, error) {
	id, err := uuid.ParseBytes(input)
	if err != nil {
		return nil, err
	}

	return FromID(id), nil
}

func MustParse(input string) UUID {
	parsed, err := Parse(input)
	if err != nil {
		return nil
	}

	return parsed
}

func MustParseBytes(input []byte) UUID {
	parsed, err := ParseBytes(input)
	if err != nil {
		return nil
	}

	return parsed
}

func String() string {
	return New().String()
}

func (u uid) Raw() uuid.UUID {
	return u.id
}

func (u uid) String() string {
	var buffer [36]byte
	encodeHex(buffer[:], u.id)
	return to.BytesToString(buffer[:])
}

func (u uid) Bytes() []byte {
	var buffer [36]byte
	encodeHex(buffer[:], u.id)
	return buffer[:]
}

func encodeHex(dst []byte, id uuid.UUID) {
	hex.Encode(dst, id[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], id[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], id[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], id[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], id[10:])
}
