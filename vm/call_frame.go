package vm

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/EclesioMeloJunior/wasvm/leb128"
)

var (
	ErrParamOutOfBounds = errors.New("param out of bounds")
	ErrWrongType        = errors.New("wrong type")
)

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
		case localGet:
			// push the parameter onto the stack.
			// advance the pointer counter to get the variable index
			c.pc += 1
			bytesRead, paramAt, err := leb128.DecodeUint(bytes.NewReader(c.instructions[c.pc:]))
			if err != nil {
				return nil, fmt.Errorf("failed to decode u32 local index: %w", err)
			}

			if len(params) < int(paramAt) {
				return nil, ErrParamOutOfBounds
			}

			c.stack.push(StackValue{
				value:   params[paramAt],
				startAt: c.pc,
				endAt:   c.pc + uint(bytesRead),
			})

			c.pc += uint(bytesRead)

		case i32Const:
			// push the i32 leb128 encoded value onto the stack.
			// lets start read the encoded number
			c.pc += 1
			bytesRead, value, err := leb128.DecodeInt[int32](
				bytes.NewReader(c.instructions[c.pc:]))

			if err != nil {
				return nil, fmt.Errorf("failed to decode int32: %w", err)
			}

			stackBasedValue := StackValue{
				value:   value,
				startAt: c.pc,
				endAt:   c.pc + uint(bytesRead),
			}

			c.stack.push(stackBasedValue)
			c.pc += uint(bytesRead)

		case i32Add:
			rhs, err := popEnsureType[int32](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			lhs, err := popEnsureType[int32](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			c.stack.push(StackValue{
				value: lhs + rhs,
			})

			c.pc++

		case i32Sub:
			rhs, err := popEnsureType[int32](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			lhs, err := popEnsureType[int32](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			c.stack.push(StackValue{
				value: lhs - rhs,
			})

			c.pc++
		case i32Mul:
			rhs, err := popEnsureType[int32](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			lhs, err := popEnsureType[int32](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			c.stack.push(StackValue{
				value: lhs * rhs,
			})

			c.pc++

		case i32LowerThanSigned:
			rhs, err := popEnsureType[int32](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			lhs, err := popEnsureType[int32](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			if rhs < lhs {
				c.stack.push(TrueStackValue)
			} else {
				c.stack.push(FalseStackValue)
			}

			c.pc++
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
