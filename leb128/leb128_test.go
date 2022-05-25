package leb128_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/EclesioMeloJunior/wasvm/leb128"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeInt(t *testing.T) {
	tests := []struct {
		input    int
		expected []byte
	}{
		{input: -math.MaxInt32, expected: []byte{0x81, 0x80, 0x80, 0x80, 0x78}},
		{input: -165675008, expected: []byte{0x80, 0x80, 0x80, 0xb1, 0x7f}},
		{input: -624485, expected: []byte{0x9b, 0xf1, 0x59}},
		{input: -16256, expected: []byte{0x80, 0x81, 0x7f}},
		{input: -4, expected: []byte{0x7c}},
		{input: -1, expected: []byte{0x7f}},
		{input: 0, expected: []byte{0x00}},
		{input: 1, expected: []byte{0x01}},
		{input: 4, expected: []byte{0x04}},
		{input: 16256, expected: []byte{0x80, 0xff, 0x0}},
		{input: 624485, expected: []byte{0xe5, 0x8e, 0x26}},
		{input: 165675008, expected: []byte{0x80, 0x80, 0x80, 0xcf, 0x0}},
		{input: math.MaxInt32, expected: []byte{0xff, 0xff, 0xff, 0xff, 0x7}},
		{input: math.MaxInt64, expected: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0}},
	}

	for _, tt := range tests {
		enc := leb128.EncodeInt(tt.input)
		require.Equal(t, tt.expected, enc)
	}
}

func TestEncodeUint(t *testing.T) {
	tests := []struct {
		input    uint
		expected []byte
	}{
		{input: 0, expected: []byte{0x00}},
		{input: 1, expected: []byte{0x01}},
		{input: 4, expected: []byte{0x04}},
		{input: 16256, expected: []byte{0x80, 0x7f}},
		{input: 624485, expected: []byte{0xe5, 0x8e, 0x26}},
		{input: 165675008, expected: []byte{0x80, 0x80, 0x80, 0x4f}},
		{input: math.MaxUint32, expected: []byte{0xff, 0xff, 0xff, 0xff, 0xf}},
		{input: math.MaxUint64, expected: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1}},
	}

	for _, tt := range tests {
		enc := leb128.EncodeUint(tt.input)
		require.Equal(t, tt.expected, enc)
	}
}

func TestDecodeUint(t *testing.T) {
	tests := []struct {
		enc      []byte
		expected uint
		wantErr  error
	}{
		{enc: []byte{0x00}, expected: 0},
		{enc: []byte{0x04}, expected: 4},
		{enc: []byte{0x01}, expected: 1},
		{enc: []byte{0x80, 0}, expected: 0},
		{enc: []byte{0x80, 0x7f}, expected: 16256},
		{enc: []byte{0xe5, 0x8e, 0x26}, expected: 624485},
		{enc: []byte{0x80, 0x80, 0x80, 0x4f}, expected: 165675008},
		{enc: []byte{0xff, 0xff, 0xff, 0xff, 0xf}, expected: 0xffffffff},
		{enc: []byte{0xff, 0xff, 0xff, 0xff, 0xf}, expected: math.MaxUint32},
		{enc: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1}, expected: math.MaxUint64},
	}

	for _, tt := range tests {
		result, err := leb128.DecodeUint(bytes.NewReader(tt.enc))
		if tt.wantErr != nil {
			assert.EqualError(t, err, tt.wantErr.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestDecodeInt32(t *testing.T) {
	tests := []struct {
		enc      []byte
		expected int32
		wantErr  error
	}{
		{enc: []byte{0x13}, expected: 19},
		{enc: []byte{0x00}, expected: 0},
		{enc: []byte{0x04}, expected: 4},
		{enc: []byte{0xFF, 0x00}, expected: 127},
		{enc: []byte{0x81, 0x01}, expected: 129},
		{enc: []byte{0x7f}, expected: -1},
		{enc: []byte{0x81, 0x7f}, expected: -127},
		{enc: []byte{0xFF, 0x7e}, expected: -129},
	}

	for _, tt := range tests {
		result, err := leb128.DecodeInt[int32](tt.enc)
		if tt.wantErr != nil {
			assert.EqualError(t, err, tt.wantErr.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestDecodeInt64(t *testing.T) {
	tests := []struct {
		enc      []byte
		expected int64
		wantErr  error
	}{
		{enc: []byte{0x00}, expected: 0},
		{enc: []byte{0x04}, expected: 4},
		{enc: []byte{0xFF, 0x00}, expected: 127},
		{enc: []byte{0x81, 0x01}, expected: 129},
		{enc: []byte{0x7f}, expected: -1},
		{enc: []byte{0x81, 0x7f}, expected: -127},
		{enc: []byte{0xFF, 0x7e}, expected: -129},
		{enc: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x7f},
			expected: -9223372036854775808},
	}

	for _, tt := range tests {
		result, err := leb128.DecodeInt[int64](tt.enc)
		if tt.wantErr != nil {
			assert.EqualError(t, err, tt.wantErr.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}
