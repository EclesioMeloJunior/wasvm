package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/EclesioMeloJunior/wasvm/leb128"
)

const (
	TypeSection     byte = 0x01
	FunctionSection byte = 0x03
	ExportSection   byte = 0x07
	CodeSection     byte = 0x0A
)

var (
	ErrBytesLen = errors.New("unexpected bytes len")
)

type SectionParser interface {
	Parse(*bytes.Reader) error
}

type BinaryReader interface {
	io.Reader
	io.ByteReader
}

type Module struct {
	Magic   uint32 // magic number >`\0asm`
	Version uint32 // version
}

type BinaryParser struct {
	filepath string
	reader   BinaryReader

	Module *Module

	Parsers map[byte]SectionParser
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

		Parsers: map[byte]SectionParser{
			TypeSection:     &TypeSectionParser{},
			FunctionSection: &FunctionSectionParser{},
			ExportSection:   &ExportSectionParser{},
			CodeSection:     &CodeSectionParser{},
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
	if errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return fmt.Errorf("cannot read section byte: %w", err)
	}

	sectionsLen, err := leb128.DecodeUint(bp.reader.(*bytes.Reader))
	if err != nil {
		return fmt.Errorf("cannot read section len: %w", err)
	}

	if sectionsLen > 0 {
		err := bp.parseSectionContents(sectionByte, sectionsLen)
		if err != nil {
			return fmt.Errorf(
				"cannot parse: %w", err)
		}
	}

	return bp.ParseSection()
}

func (bp *BinaryParser) parseSectionContents(sectionID byte, sectionLen uint) error {
	contents := make([]byte, sectionLen)
	n, err := bp.reader.Read(contents)

	if err != nil {
		return fmt.Errorf("cannot read section contents: %w", err)
	} else if n != int(sectionLen) {
		return fmt.Errorf("expected %d bytes. read %d bytes", sectionLen, n)
	}

	parser, ok := bp.Parsers[sectionID]
	if !ok {
		return fmt.Errorf("empty parser for section ID 0x%x", sectionID)
	}

	err = parser.Parse(bytes.NewReader(contents))
	if err != nil {
		return fmt.Errorf("failed while parsing section 0x%x: %w", sectionID, err)
	}

	bp.Parsers[sectionID] = parser
	return nil
}
