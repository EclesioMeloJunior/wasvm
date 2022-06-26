package parser_test

import (
	"encoding/binary"
	"testing"

	"github.com/EclesioMeloJunior/wasvm/parser"

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

func TestSimpleWasm_BinaryParse_Sections(t *testing.T) {
	wasm, err := parser.NewBinaryParser(simpleWasm)
	assert.NoError(t, err)

	err = wasm.ParseMagicNumber()
	assert.NoError(t, err)

	err = wasm.ParseVersion()
	assert.NoError(t, err)

	err = wasm.ParseSection()
	assert.NoError(t, err)

	typeSection := wasm.Parsers[parser.TypeSection].(*parser.TypeSectionParser)
	assert.Len(t, typeSection.Types, 1)

	for _, ttype := range typeSection.Types {
		fs, ok := ttype.(*parser.FunctionSignatureParser)
		assert.True(t, ok)

		assert.Len(t, fs.ParamsTypes, 0)
		assert.Len(t, fs.ResultsTypes, 1)

		rt := fs.ResultsTypes[0]
		assert.Equal(t, rt.SpecByte, parser.I32_NUM_TYPE)
		assert.Equal(t, rt.SpecType, parser.NumType)
	}

	functionSection := wasm.Parsers[parser.FunctionSection].(*parser.FunctionSectionParser)
	assert.Len(t, functionSection.Funcs, 1)

	for _, f := range functionSection.Funcs {
		assert.Equal(t, f.TypeIndex, 0)
	}

	exportSection := wasm.Parsers[parser.ExportSection].(*parser.ExportSectionParser)
	assert.Len(t, exportSection.Exports, 1)

	for _, exported := range exportSection.Exports {
		assert.Equal(t, exported.Name, "helloWorld")
		assert.Equal(t, exported.Index, 0)
		assert.Equal(t, exported.Type, parser.ExportedFunc)
	}

	codeSection := wasm.Parsers[parser.CodeSection].(*parser.CodeSectionParser)
	assert.Len(t, codeSection.FunctionsCode, 1)

	for _, code := range codeSection.FunctionsCode {
		assert.Len(t, code.Locals, 0)
		assert.Len(t, code.Body, 3)

		expectedBody := []byte{0x41, 0x2A, 0x0B}
		for idx, instr := range expectedBody {
			assert.Equal(t, code.Body[idx], instr)
		}
	}
}
