package gowasm

import (
	"os"

	"github.com/EclesioMeloJunior/gowasm/binary"
)

type WasmFormat byte

const (
	BinaryFormat WasmFormat = iota
	TextFormat
)

func Load(path string, ttype WasmFormat) (contents []byte, err error) {
	switch ttype {
	case BinaryFormat:
		binary.Parse(path)
	case TextFormat:
		binary.Parse(path)
	}

	contents, err = os.ReadFile(path)
	return
}
