package parser

import (
	"errors"
	"fmt"
)

var (
	ErrFunctionWithouSignature = errors.New("function does not have a respective signature")
	ErrFunctionWithouCode      = errors.New("function does not have a respective code")
)

func BinaryFormat(filepath string) (*BinaryParser, error) {
	bp, err := NewBinaryParser(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot instantiate a binary parser: %w", err)
	}

	// starting parsing the `wasm header` values
	if err := bp.ParseMagicNumber(); err != nil {
		return nil, fmt.Errorf("cannot parse magic number: %w", err)
	}

	if err := bp.ParseVersion(); err != nil {
		return nil, fmt.Errorf("cannot parse version: %w", err)
	}

	if err := bp.ParseSection(); err != nil {
		return nil, fmt.Errorf("cannot parse section: %w", err)
	}

	if err := bondFunctionSignatureAndCode(bp); err != nil {
		return nil, fmt.Errorf("cannot bond function parts: %w", err)
	}

	return bp, nil
}

func bondFunctionSignatureAndCode(bp *BinaryParser) error {
	functionSection := bp.Parsers[FunctionSection].(*FunctionSectionParser)
	if len(functionSection.Funcs) < 1 {
		return nil
	}

	typeSection := bp.Parsers[TypeSection].(*TypeSectionParser)
	codeSection := bp.Parsers[CodeSection].(*CodeSectionParser)

	for idx, function := range functionSection.Funcs {
		if len(typeSection.Types) < function.TypeIndex {
			return fmt.Errorf("%w: %d", ErrFunctionWithouSignature, function.TypeIndex)
		}

		ttype := typeSection.Types[function.TypeIndex]
		signature, ok := ttype.(*FunctionSignatureParser)
		if !ok {
			return fmt.Errorf("%w: expected *FunctionSignatureParser, got: %T",
				ErrFunctionWithouSignature, ttype)
		}

		// each index correspond to a code in the code section, if there is no
		// code for the current index then we must return an error
		if len(codeSection.FunctionsCode) < idx {
			return fmt.Errorf("%w: %d", ErrFunctionWithouCode, function.TypeIndex)
		}

		code := codeSection.FunctionsCode[idx]

		functionSection.Funcs[idx].Code = code
		functionSection.Funcs[idx].Signature = signature
	}

	return nil
}
