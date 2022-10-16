package opcodes

import (
	"fmt"
)

type OpCode byte

func (i OpCode) String() string {
	switch i {
	case I32Const:
		return "i32.const"
	case I32Add:
		return "i32.add"
	case I32Sub:
		return "i32.sub"
	case I32Mul:
		return "i32.mul"
	case I32LowerThanSigned:
		return "i32.lt_s"
	case If:
		return "if"
	case Else:
		return "else"
	case End:
		return "end"
	case Call:
		return "call"
	case Return:
		return "return"
	case EmptyBlockType:
		return "empty(blocktype)"
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

	If     OpCode = 0x04
	Else   OpCode = 0x05
	End    OpCode = 0x0B
	Return OpCode = 0x0F

	Call OpCode = 0x10

	EmptyBlockType = 0x40
)
