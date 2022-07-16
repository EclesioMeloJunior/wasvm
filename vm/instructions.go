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
	localGet           Instruction = 0x20
	i32Const           Instruction = 0x41
	i32Add             Instruction = 0x6A
	i32Sub             Instruction = 0x6B
	i32Mul             Instruction = 0x6C
	i32LowerThanSigned Instruction = 0x48
	End                Instruction = 0x0B
)
