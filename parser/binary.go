package parser

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

var (
	ErrBytesLen = errors.New("unexpected bytes len")
)

type Wasm struct {
	filepath string
	reader   *bufio.Reader

	Module *Module
}

func NewWasm(filepath string) (*Wasm, error) {
	fbytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	return &Wasm{
		Module:   &Module{},
		filepath: filepath,
		reader:   bufio.NewReader(bytes.NewReader(fbytes)),
	}, nil
}

func (w *Wasm) ParseMagicNumber() error {
	const magicNumberLen = 4
	magicBytes := make([]byte, magicNumberLen)

	n, err := w.reader.Read(magicBytes)
	if err != nil {
		return fmt.Errorf("cannot read magic number: %w", err)
	}

	if n != magicNumberLen {
		return fmt.Errorf("%w: %d, expected: %d",
			ErrBytesLen, n, magicNumberLen)
	}

	w.Module.Magic = binary.LittleEndian.Uint32(magicBytes)
	return nil
}

func (w *Wasm) ParseVersion() error {
	const versionBytesLen = 4
	versionBytes := make([]byte, versionBytesLen)

	n, err := w.reader.Read(versionBytes)
	if err != nil {
		return fmt.Errorf("cannot read version: %w", err)
	}

	if n != versionBytesLen {
		return fmt.Errorf("%w : %d, expected: %d",
			ErrBytesLen, n, versionBytesLen)
	}

	w.Module.Version = binary.LittleEndian.Uint32(versionBytes)
	return nil
}

func (w *Wasm) Next() (b byte, err error) {
	b, err = w.reader.ReadByte()
	if err != nil {
		return b, fmt.Errorf("cannot read next wasm byte: %w", err)
	}

	return b, nil
}
