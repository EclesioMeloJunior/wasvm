package leb128

type Encoded byte

const (
	Signed Encoded = iota
	Unsigned
)

func Encode(n []byte, ttype Encoded) {

}

func encodeUnsigned()
