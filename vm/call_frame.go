package vm

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/EclesioMeloJunior/wasvm/leb128"
)

var ErrStackOverflow = errors.New("stackoverflow")
var ErrEmptyStack = errors.New("empty stack")

type StackValue struct {
	value   any
	startAt uint
	endAt   uint
}

type Stack []StackValue

func (s *Stack) push(value StackValue) error {
	defPointer := *s
	if len(defPointer) >= cap(defPointer) {
		return fmt.Errorf("%w: limit %d", ErrStackOverflow, cap(defPointer))
	}

	*s = append(defPointer, value)
	return nil
}

func (s *Stack) pop() (value StackValue, err error) {
	defPointer := *s
	if len(defPointer) == 0 {
		return value, ErrEmptyStack
	}

	removeAt := len(defPointer) - 1
	value = defPointer[removeAt]

	*s = defPointer[:removeAt]
	return value, nil
}

type callFrame struct {
	pc    uint
	stack Stack

	params       []any
	results      []any
	instructions []byte
}

func (c *callFrame) Call(params ...any) ([]any, error) {
	for {
		currentInstruction := Instruction(c.instructions[c.pc])

		switch currentInstruction {
		case i32Const:
			// lets start read the encoded number
			c.pc += 1
			bytesRead, value, err := leb128.DecodeInt[int32](
				bytes.NewReader(c.instructions[c.pc:]))

			if err != nil {
				return nil, fmt.Errorf("failed to decode int32: %w", err)
			}

			stackBasedValue := StackValue{
				value:   value,
				startAt: uint(c.pc),
				endAt:   uint(c.pc + uint(bytesRead)),
			}

			c.stack.push(stackBasedValue)
			c.pc += uint(bytesRead)
		case End:
			if len(c.results) > 0 && len(c.stack) == 0 {
				return nil, fmt.Errorf("stack empty but expected %d return(s)",
					len(c.results))
			}

			results := make([]any, len(c.results))
			for idx := 0; idx < len(c.results); idx++ {
				popped, err := c.stack.pop()
				if err != nil {
					return nil, fmt.Errorf("cannot pop result from stack: %w", err)
				}

				results[idx] = popped.value
			}

			return results, nil

		default:
			return nil, fmt.Errorf("unknonw instruction: %s", currentInstruction)
		}
	}
}
