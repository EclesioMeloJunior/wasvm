package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/EclesioMeloJunior/wasvm/leb128"
)

type FunctionSignatureParser struct {
	ParamsTypes  []Type
	ResultsTypes []Type
}

func (f FunctionSignatureParser) String() string {
	result := "func("
	for idx, p := range f.ParamsTypes {
		result += p.String()
		if idx < len(f.ParamsTypes)-1 {
			result += ", "
		}
	}
	result += ")"

	if len(f.ResultsTypes) > 0 {
		result += " -> ("
	}

	for idx, p := range f.ResultsTypes {
		result += p.String()
		if idx < len(f.ParamsTypes)-1 {
			result += ", "
		}
	}

	result += ")"

	return result
}

func (f *FunctionSignatureParser) Parse(b *bytes.Reader) error {
	_, paramsLen, err := leb128.DecodeUint(b)
	if err != nil {
		return fmt.Errorf("cannot read params length: %w", err)
	}

	paramsTypes := make([]Type, paramsLen)
	for i := 0; i < int(paramsLen); i++ {
		paramType, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("cannot read param type at %d: %w", i, err)
		}

		switch paramType {
		case I32_NUM_TYPE, I64_NUM_TYPE, F32_NUM_TYPE, F64_NUM_TYPE:
			paramsTypes[i] = Type{
				SpecType: NumType,
				SpecByte: paramType,
			}
		}
	}

	_, resultsLen, err := leb128.DecodeUint(b)
	if err != nil {
		return fmt.Errorf("cannot read results length: %w", err)
	}

	resultsTypes := make([]Type, resultsLen)
	for i := 0; i < int(resultsLen); i++ {
		resultType, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("cannot read result type at %d: %w", i, err)
		}

		switch resultType {
		case I32_NUM_TYPE, I64_NUM_TYPE, F32_NUM_TYPE, F64_NUM_TYPE:
			resultsTypes[i] = Type{
				SpecType: NumType,
				SpecByte: resultType,
			}
		}
	}

	f.ParamsTypes = paramsTypes
	f.ResultsTypes = resultsTypes

	return nil
}

type TypeSectionParser struct {
	Types []Parser
}

func (t *TypeSectionParser) Parse(b *bytes.Reader) error {
	_, typeSectionLen, err := leb128.DecodeUint(b)
	if err != nil {
		return fmt.Errorf("cannot read type section length: %w", err)
	}

	funcSignatureTypes := make([]Parser, 0, typeSectionLen)

	for i := 0; i < int(typeSectionLen); i++ {
		typeTag, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("failed to read tag at index %d: %w", i, err)
		}

		if typeTag == FunctionTag {
			functionSigParser := &FunctionSignatureParser{}
			if err := functionSigParser.Parse(b); err != nil {
				return fmt.Errorf("cannot parse function signature at index %d: %w", i, err)
			}

			funcSignatureTypes = append(funcSignatureTypes, functionSigParser)
		}
	}

	t.Types = funcSignatureTypes
	return nil
}

type Function struct {
	TypeIndex int
	Signature *FunctionSignatureParser
	Code      *CodeParser
}

type FunctionSectionParser struct {
	Funcs []*Function
}

func (f *FunctionSectionParser) Parse(b *bytes.Reader) error {
	_, funcsLen, err := leb128.DecodeUint(b)
	if err != nil {
		return fmt.Errorf("cannot read function amount: %w", err)
	}

	funcs := make([]*Function, funcsLen)
	for i := 0; i < int(funcsLen); i++ {
		_, typeIndex, err := leb128.DecodeUint(b)
		if err != nil {
			return fmt.Errorf("cannot read function type index at %d: %w", i, err)
		}

		funcs[i] = &Function{
			TypeIndex: int(typeIndex),
		}
	}

	f.Funcs = funcs
	return nil
}

// Type tells us what is being exported
// 0x00 funcidx
// 0x01 tableidx
// 0x02 memidx
// 0x03 globalidx
type ExportedType byte

const (
	ExportedFunc   ExportedType = 0x00
	ExportedTable  ExportedType = 0x01
	ExportedMem    ExportedType = 0x02
	ExportedGlobal ExportedType = 0x03
)

type Export struct {
	Name  string
	Type  ExportedType
	Index int
}

type ExportSectionParser struct {
	Exports []*Export
}

func (e *ExportSectionParser) Parse(b *bytes.Reader) error {
	_, exportsLen, err := leb128.DecodeUint(b)
	if err != nil {
		return fmt.Errorf("cannot read number of exports: %w", err)
	}

	exports := make([]*Export, exportsLen)

	for i := 0; i < int(exportsLen); i++ {
		_, nameLen, err := leb128.DecodeUint(b)
		if err != nil {
			return fmt.Errorf("cannot read exported name length at %d: %w", i, err)
		}

		nameBytes := make([]byte, nameLen)
		n, err := b.Read(nameBytes)
		if err != nil {
			return fmt.Errorf("cannot read exported name bytes at %d: %w", i, err)
		} else if n != int(nameLen) {
			return fmt.Errorf("expected name bytes length %d. got %d", nameLen, n)
		}

		exportType, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("cannot read exported type at %d: %w", i, err)
		}

		_, exportIdx, err := leb128.DecodeUint(b)
		if err != nil {
			return fmt.Errorf("cannot read exported index at %d: %w", i, err)
		}

		exports[i] = &Export{
			Index: int(exportIdx),
			Type:  ExportedType(exportType),
			Name:  string(nameBytes),
		}
	}

	e.Exports = exports
	return nil
}

type CodeParser struct {
	Body   []byte
	Locals []Type
}

func (c *CodeParser) Parse(b *bytes.Reader, len uint) error {
	_, localsLen, err := leb128.DecodeUint(b)
	if err != nil {
		return fmt.Errorf("cannot read local length: %w", err)
	}

	if localsLen > 0 {
		c.parseLocals(b, localsLen)
	}

	body := make([]byte, 0)
	for i := 0; i < int(len); i++ {
		b, err := b.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		body = append(body, b)
	}

	c.Body = body
	return nil
}

func (c *CodeParser) parseLocals(b *bytes.Reader, len uint) error {
	c.Locals = make([]Type, len)

	for i := uint(0); i < len; i++ {
		localType, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("while reading local type byte: %w", err)
		}

		switch localType {
		case I32_NUM_TYPE, I64_NUM_TYPE:
			c.Locals[i] = Type{
				SpecType: NumType,
				SpecByte: localType,
			}
		default:
			unsupportedType := Type{
				SpecByte: localType,
			}

			panic(fmt.Sprintf("local type not supported yet: %s (%v)",
				unsupportedType.String(), localType))
		}
	}

	return nil
}

type CodeSectionParser struct {
	FunctionsCode []*CodeParser
}

func (c *CodeSectionParser) Parse(b *bytes.Reader) error {
	_, amount, err := leb128.DecodeUint(b)
	if err != nil {
		return fmt.Errorf("cannot read number of functions: %w", err)
	}

	codes := make([]*CodeParser, amount)
	for i := 0; i < int(amount); i++ {
		_, totalCodeSize, err := leb128.DecodeUint(b)
		if err != nil {
			return fmt.Errorf("cannot read the code length at %d: %w", i, err)
		}

		code := make([]byte, totalCodeSize)
		n, err := b.Read(code)
		if err != nil {
			return fmt.Errorf("cannot read the code at %d: %w", i, err)
		} else if n != int(totalCodeSize) {
			return fmt.Errorf("expected code bytes length %d. got %d", totalCodeSize, n)
		}

		codeParser := &CodeParser{}
		err = codeParser.Parse(bytes.NewReader(code), totalCodeSize)
		if err != nil {
			return fmt.Errorf("cannot parse code instructions at %d: %w", i, err)
		}

		codes[i] = codeParser
	}

	c.FunctionsCode = codes

	return nil
}

type StartSectionParser struct{}

func (i *StartSectionParser) Parse(b *bytes.Reader) error {
	return nil
}

type ImportsSectionParser struct{}

func (i *ImportsSectionParser) Parse(b *bytes.Reader) error {
	return nil
}
