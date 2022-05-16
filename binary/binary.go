package binary

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type Wasm struct {
	filepath string
	reader   *bufio.Reader
}

func NewWasm(filepath string) (*Wasm, error) {
	fbytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	return &Wasm{
		filepath: filepath,
		reader:   bufio.NewReader(bytes.NewReader(fbytes)),
	}, nil
}

func (w *Wasm) Next() (b byte, err error) {
	b, err = w.reader.ReadByte()
	if err != nil {
		return b, fmt.Errorf("cannot read next wasm byte: %w", err)
	}

	return b, nil
}
