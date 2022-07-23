package vm

import (
	"errors"
	"fmt"

	"github.com/EclesioMeloJunior/wasvm/parser"
)

var ErrCannotExportFunction = errors.New("cannot export function")

type exportedFunction func(...any) any

type Runtime struct {
	binary   *parser.BinaryParser
	Exported map[string]*callFrame
}

func NewRuntime(bp *parser.BinaryParser) (*Runtime, error) {
	runtime := &Runtime{
		binary: bp,
	}

	if err := exposeExportedFunctions(runtime); err != nil {
		return nil, err
	}

	return runtime, nil
}

func exposeExportedFunctions(runtime *Runtime) error {
	functionSection := runtime.binary.Parsers[parser.FunctionSection].(*parser.FunctionSectionParser)
	codeSection := runtime.binary.Parsers[parser.CodeSection].(*parser.CodeSectionParser)
	exportedSection := runtime.binary.Parsers[parser.ExportSection].(*parser.ExportSectionParser)

	runtime.Exported = make(map[string]*callFrame, len(exportedSection.Exports))

	for _, exported := range exportedSection.Exports {
		switch exported.Type {
		case parser.ExportedFunc:
			if len(codeSection.FunctionsCode) < exported.Index {
				return fmt.Errorf("%w: expected function on index %d",
					ErrCannotExportFunction, exported.Index)
			}

			exportedFunction := functionSection.Funcs[exported.Index]
			exportedCode := codeSection.FunctionsCode[exported.Index]

			runtime.Exported[exported.Name] = newCallFrame(exportedCode.Body,
				exportedFunction.Signature.ParamsTypes,
				exportedFunction.Signature.ResultsTypes)
		}
	}

	return nil
}
