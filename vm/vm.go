package vm

import (
	"errors"
	"fmt"

	"github.com/EclesioMeloJunior/wasvm/parser"
)

type exportedFunction func(...any) any

var (
	ErrCannotExportFunction = errors.New("cannot export function")
)

type callFrame struct {
	stack []byte
	pc    uint

	params       []any
	results      []any
	instructions []byte
}

func (c *callFrame) Call(params ...any) {
	for _, instructionByte := range c.instructions {
		fmt.Printf("%x\n", instructionByte)
	}
}

type Runtime struct {
	binary   *parser.BinaryParser
	Exported map[string]*callFrame
}

func NewRuntime(bp *parser.BinaryParser) (*Runtime, error) {
	if err := bp.ParseMagicNumber(); err != nil {
		return nil, err
	}

	if err := bp.ParseVersion(); err != nil {
		return nil, err
	}

	if err := bp.ParseVersion(); err != nil {
		return nil, err
	}

	runtime := &Runtime{
		binary: bp,
	}

	if err := exposeExportedFunctions(runtime); err != nil {
		return nil, err
	}

	return runtime, nil
}

func (r *Runtime) engine() {}

func exposeExportedFunctions(runtime *Runtime) error {
	functionSection := runtime.binary.Parsers[parser.FunctionSection].(*parser.FunctionSectionParser)
	exportedSection := runtime.binary.Parsers[parser.ExportSection].(*parser.ExportSectionParser)

	for _, exported := range exportedSection.Exports {
		switch exported.Type {
		case parser.ExportedFunc:
			if len(functionSection.Funcs) < exported.Index {
				return fmt.Errorf("%w: expected function on index %d",
					ErrCannotExportFunction, exported.Index)
			}

			eFunc := functionSection.Funcs[exported.Index]
			runtime.Exported[exported.Name] = buildExportedFunction(eFunc)
		}
	}

	return nil
}

func buildExportedFunction(f *parser.Function) *callFrame {
	cf := &callFrame{
		pc:           0,
		stack:        make([]byte, 1024),
		instructions: f.Code.Body,
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
			cf.params[idx] = int32(0)
		case parser.I64_NUM_TYPE:
			cf.params[idx] = int64(0)
		default:
			// TODO: implement other types at
			panic(fmt.Sprintf("result type not supported yet: %x", rt.SpecByte))
		}
	}

	return nil
}
