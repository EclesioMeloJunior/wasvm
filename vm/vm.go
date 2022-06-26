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
			runtime.Exported[exported.Name] = buildExportedFunction(exportedFunction, exportedCode)
		}
	}

	return nil
}

func buildExportedFunction(f *parser.Function, c *parser.CodeParser) *callFrame {
	cf := &callFrame{
		pc:           0,
		stack:        make([]StackValue, 0, 1024),
		instructions: c.Body,
		params:       make([]any, len(f.Signature.ParamsTypes)),
		results:      make([]any, len(f.Signature.ResultsTypes)),
	}

	for idx, pt := range f.Signature.ParamsTypes {
		switch pt.SpecByte {
		case parser.I32_NUM_TYPE:
			cf.params[idx] = int32(0)
		case parser.I64_NUM_TYPE:
			cf.params[idx] = int64(0)
		default:
			// TODO: implement other types at
			panic(fmt.Sprintf("param type not supported yet: %x", pt.SpecByte))
		}
	}

	for idx, rt := range f.Signature.ResultsTypes {
		switch rt.SpecByte {
		case parser.I32_NUM_TYPE:
			cf.results[idx] = int32(0)
		case parser.I64_NUM_TYPE:
			cf.results[idx] = int64(0)
		default:
			// TODO: implement other types at
			panic(fmt.Sprintf("result type not supported yet: %x", rt.SpecByte))
		}
	}

	return cf
}
