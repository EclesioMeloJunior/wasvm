package parser_test

import (
	"encoding/binary"
	"testing"

	"github.com/EclesioMeloJunior/gowasm/parser"

	"github.com/stretchr/testify/assert"
)

const simpleWasm = "../tests/simple.wasm"

func TestBinaryParse_MagicNumber_VersionNumber(t *testing.T) {
	wasm, err := parser.NewBinaryParser(simpleWasm)
	assert.NoError(t, err)

	err = wasm.ParseMagicNumber()
	assert.NoError(t, err)

	magic := make([]byte, 4)
	binary.LittleEndian.PutUint32(magic, wasm.Module.Magic)

	assert.Equal(t, string(magic), "\x00asm")
	assert.NoError(t, err)

	err = wasm.ParseVersion()
	assert.NoError(t, err)

	assert.Equal(t, uint32(1), wasm.Module.Version)
}

func TestBinaryParse_Sections(t *testing.T) {
	wasm, err := parser.NewBinaryParser(simpleWasm)
	assert.NoError(t, err)

	err = wasm.ParseMagicNumber()
	assert.NoError(t, err)

	err = wasm.ParseVersion()
	assert.NoError(t, err)

	err = wasm.ParseSection()
	assert.NoError(t, err)
}
