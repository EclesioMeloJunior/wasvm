package vm

import (
	"bytes"
	"fmt"

	"github.com/EclesioMeloJunior/wasvm/leb128"
)

type StackValue struct {
	value   any
	startAt uint
	endAt   uint
}

type callFrame struct {
	stack []StackValue
	pc    uint

	params       []any
	results      []any
	instructions []byte
}

func (c *callFrame) Call(params ...any) (any, error) {
	for {
		currentInstruction := Instruction(c.instructions[c.pc])

		switch currentInstruction {
		case i32Const:
			bytesRead, value, err := leb128.DecodeInt[int32](
				bytes.NewReader(c.instructions[c.pc+1:]))
			if err != nil {
				return nil, fmt.Errorf("failed to decode int32: %w", err)
			}

			stackBasedValue := StackValue{
				value:   value,
				startAt: uint(c.pc + 1),
				endAt:   uint(c.pc + 1 + uint(bytesRead)),
			}

			c.stack = append(c.stack, stackBasedValue)
			c.pc += uint(bytesRead)

		case End:
			if len(c.results) > 0 && len(c.stack) == 0 {
				return nil, fmt.Errorf("stack empty but expected %d return(s)", len(c.results))
			}

		default:
			return nil, fmt.Errorf("unknonw instruction: %s (%v)",
				currentInstruction, currentInstruction)
		}
	}
}
