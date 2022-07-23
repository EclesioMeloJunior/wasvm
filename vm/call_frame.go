package vm

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/EclesioMeloJunior/wasvm/leb128"
	"github.com/EclesioMeloJunior/wasvm/opcodes"
	"github.com/EclesioMeloJunior/wasvm/parser"
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

func newCallFrame(instructions []byte, paramTypes, resultTypes []parser.Type) *callFrame {
	cf := &callFrame{
		pc:           0,
		stack:        make([]StackValue, 0, 1024),
		instructions: instructions,
		params:       make([]any, len(paramTypes)),
		results:      make([]any, len(resultTypes)),
	}

	for idx, pt := range paramTypes {
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

	for idx, rt := range resultTypes {
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

func (c *callFrame) searchInstruction(oc opcodes.OpCode) (position uint) {
	startAt := c.pc + 1

	for i, inst := range c.instructions[startAt:] {
		switch inst {
		case byte(oc):
			return startAt + uint(i)
		}
	}

	return 0
}

func (c *callFrame) Call(params ...any) ([]any, error) {
	for {
		if uint(len(c.instructions)) <= c.pc {
			return nil, nil
		}

		currentInstruction := opcodes.OpCode(c.instructions[c.pc])

		switch currentInstruction {
		case opcodes.LocalGet:
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

		case opcodes.I32Const:
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

		case opcodes.I32Add:
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

		case opcodes.I32Sub:
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

		case opcodes.I32Mul:
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

		case opcodes.I32LowerThanSigned:
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

		case opcodes.If:
			jumpToElse := c.searchInstruction(opcodes.Else)
			jumpToIfEnd := c.searchInstruction(opcodes.End)
			if jumpToIfEnd == 0 {
				return nil, fmt.Errorf("failed to find if end")
			}

			condition, err := popEnsureType[bool](&c.stack)
			if err != nil {
				return nil, fmt.Errorf("cannot pop: %w", err)
			}

			// check if the IF branch contains a result type
			resultTypeByteCode := c.instructions[c.pc+1]
			resultType := make([]parser.Type, 0, 1)

			switch resultTypeByteCode {
			case parser.I32_NUM_TYPE, parser.I64_NUM_TYPE, parser.F32_NUM_TYPE, parser.F64_NUM_TYPE:
				resultType = append(resultType, parser.Type{
					SpecType: parser.NumType,
					SpecByte: resultTypeByteCode,
				})
			}

			// compute the if branch call frame
			if condition {
				var amountOfInstructions uint
				if jumpToElse != 0 {
					amountOfInstructions = jumpToElse - c.pc
				} else {
					amountOfInstructions = jumpToIfEnd - c.pc
				}

				branchInstructions := make([]byte, amountOfInstructions-1)
				if len(resultType) > 0 {
					// if the IF branch contains a result type
					// then include the instructions without the RESULT TYPE or ELSE opcode
					copy(branchInstructions[:], c.instructions[c.pc+2:c.pc+amountOfInstructions])
				} else {
					// only removes the ELSE opcode
					copy(branchInstructions[:], c.instructions[c.pc+1:c.pc+amountOfInstructions])
				}

				latestItem := amountOfInstructions - 2
				branchInstructions[latestItem] = byte(opcodes.End)
				branchCallFrame := newCallFrame(branchInstructions, nil, resultType)

				result, err := branchCallFrame.Call(params...)
				if err != nil {
					return nil, fmt.Errorf("fail to call if branch call frame: %w", err)
				}

				// only push the result to the stack if we spec a result from if branch call frame
				if len(result) > 0 && len(resultType) > 0 {
					c.stack.push(StackValue{
						value: result[0],
					})
				}
			} else if jumpToElse != 0 { // there is a else branch
				amountOfInstructions := jumpToIfEnd - jumpToElse
				branchInstructions := make([]byte, amountOfInstructions)
				copy(branchInstructions[:], c.instructions[jumpToElse+1:jumpToElse+1+amountOfInstructions])

				branchCallFrame := newCallFrame(branchInstructions, nil, resultType)
				result, err := branchCallFrame.Call(params...)
				if err != nil {
					return nil, fmt.Errorf("fail to call if branch call frame: %w", err)
				}

				// only push the result to the stack if we spec a result from if branch call frame
				if len(result) > 0 && len(resultType) > 0 {
					c.stack.push(StackValue{
						value: result[0],
					})
				}
			}

			c.pc = jumpToIfEnd + 1
		case opcodes.End:
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
