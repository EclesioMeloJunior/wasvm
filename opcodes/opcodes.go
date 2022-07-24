package opcodes

import "fmt"

type OpCode byte

func (i OpCode) String() string {
	switch i {
	case I32Const:
		return "i32.const"
	default:
		return fmt.Sprintf("%x", byte(i))
	}
}

const (
	LocalGet OpCode = 0x20

	I32Const           OpCode = 0x41
	I32Add             OpCode = 0x6A
	I32Sub             OpCode = 0x6B
	I32Mul             OpCode = 0x6C
	I32LowerThanSigned OpCode = 0x48

	If   OpCode = 0x04
	Else OpCode = 0x05
	End  OpCode = 0x0B

	Call OpCode = 0x10
)
