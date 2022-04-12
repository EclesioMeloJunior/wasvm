package leb128_test

import (
	"encoding/binary"
	"testing"

	"github.com/EclesioMeloJunior/gowasm/leb128"
)

func TestUnsignedLEB128(t *testing.T) {
	tests := []struct {
		number   uint32
		expected []byte
	}{
		{624485, []byte{0xE5, 0x8E, 0x26}},
	}

	for _, tt := range tests {
		var 
		binary.BigEndian.PutUint32(, tt.number)
		leb128.Encode()
	}
}
