package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/EclesioMeloJunior/gowasm/leb128"
)

type SectionParser interface {
	Parse(*bytes.Reader) error
}

type BinaryReader interface {
	io.Reader
	io.ByteReader
}

var (
	ErrBytesLen = errors.New("unexpected bytes len")
)

type Module struct {
	Magic   uint32 // magic number >`\0asm`
	Version uint32 // version
}

type BinaryParser struct {
	filepath string
	reader   BinaryReader

	Module *Module

	sectionParsers map[byte]SectionParser
}

func NewBinaryParser(filepath string) (*BinaryParser, error) {
	fbytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	return &BinaryParser{
		Module:   &Module{},
		filepath: filepath,
		reader:   bytes.NewReader(fbytes),

		sectionParsers: map[byte]SectionParser{
			0x01: &TypeSectionParser{},
		},
	}, nil
}

func (bp *BinaryParser) ParseMagicNumber() error {
	const magicNumberLen = 4
	magicBytes := make([]byte, magicNumberLen)

	n, err := bp.reader.Read(magicBytes)
	if err != nil {
		return fmt.Errorf("cannot read magic number: %w", err)
	}

	if n != magicNumberLen {
		return fmt.Errorf("%w: %d, expected: %d",
			ErrBytesLen, n, magicNumberLen)
	}

	bp.Module.Magic = binary.LittleEndian.Uint32(magicBytes)
	return nil
}

func (bp *BinaryParser) ParseVersion() error {
	const versionBytesLen = 4
	versionBytes := make([]byte, versionBytesLen)

	n, err := bp.reader.Read(versionBytes)
	if err != nil {
		return fmt.Errorf("cannot read version: %w", err)
	}

	if n != versionBytesLen {
		return fmt.Errorf("%w : %d, expected: %d",
			ErrBytesLen, n, versionBytesLen)
	}

	bp.Module.Version = binary.LittleEndian.Uint32(versionBytes)
	return nil
}

func (bp *BinaryParser) ParseSection() error {
	sectionByte, err := bp.reader.ReadByte()
	if err != nil {
		return fmt.Errorf("cannot read section byte: %w", err)
	}

	contentsLen, err := leb128.DecodeUint(bp.reader.(*bytes.Reader))
	if err != nil {
		return fmt.Errorf("cannot read section len: %w", err)
	}

	if contentsLen > 0 {
		bp.parseSectionContents(sectionByte, contentsLen)
	}

	return nil
}

func (bp *BinaryParser) parseSectionContents(sectionID byte, sectionLen uint) error {
	contents := make([]byte, sectionLen)
	n, err := bp.reader.Read(contents)

	if err != nil {
		return fmt.Errorf("cannot read section contents: %w", err)
	} else if n != int(sectionLen) {
		return fmt.Errorf("expected %d bytes. read %d bytes", sectionLen, n)
	}

	parser, ok := bp.sectionParsers[sectionID]
	if !ok {
		return fmt.Errorf("empty parser for section ID %x", sectionID)
	}

	parser.Parse(bytes.NewReader(contents))
	return nil
}
