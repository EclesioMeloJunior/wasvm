package vm

import "fmt"

type Instruction byte

func (i Instruction) String() string {
	switch i {
	case i32Const:
		return "i32.const"
	default:
		return fmt.Sprintf("%x", byte(i))
	}
}

const (
	i32Const Instruction = 0x41
	End      Instruction = 0x0B
)
