package parser

import (
	"bytes"
	"fmt"

	"github.com/EclesioMeloJunior/gowasm/leb128"
)

type FunctionSignatureParser struct {
	Tag          byte
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

func (f *FunctionSignatureParser) Parser(b *bytes.Reader) error {
	paramsLen, err := leb128.DecodeUint(b)
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

	resultsLen, err := leb128.DecodeUint(b)
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
	Types []*FunctionSignatureParser
}

func (t *TypeSectionParser) Parse(b *bytes.Reader) error {
	typeSectionLen, err := leb128.DecodeUint(b)
	if err != nil {
		return fmt.Errorf("cannot read type section length: %w", err)
	}

	functions := make([]*FunctionSignatureParser, 0, typeSectionLen)

	for i := 0; i < int(typeSectionLen); i++ {
		typeTag, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("failed to read tag at index %d: %w", i, err)
		}

		if typeTag == FunctionTag {
			functionSigParser := &FunctionSignatureParser{
				Tag: FunctionTag,
			}

			if err := functionSigParser.Parser(b); err != nil {
				return fmt.Errorf("cannot parse function signature at index %d: %w", i, err)
			}

			functions = append(functions, functionSigParser)
		}
	}

	for _, f := range functions {
		fmt.Printf("%s\n", f)
	}

	return nil
}
