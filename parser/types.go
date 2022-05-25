package parser

type ValType byte

const (
	NumType ValType = iota
	VecType
	RefType
)

const (
	I32_NUM_TYPE    byte = 0x7F
	I64_NUM_TYPE         = 0x7E
	F32_NUM_TYPE         = 0x7D
	F64_NUM_TYPE         = 0x7C
	VEC_TYPE             = 0x7B
	FUNC_REF_TYPE        = 0x70
	EXTERN_REF_TYPE      = 0x6F

	FunctionTag = 0x60
)

type Type struct {
	SpecType ValType
	SpecByte byte
}

func (v Type) String() string {
	switch v.SpecByte {
	case I32_NUM_TYPE:
		return "i32"
	case I64_NUM_TYPE:
		return "i64"
	case F32_NUM_TYPE:
		return "f32"
	case F64_NUM_TYPE:
		return "f64"
	}

	return "?"
}
