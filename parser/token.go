package parser

const (
	I32_NUM_TYPE = 0x7F
	I64_NUM_TYPE = 0x7E
	F32_NUM_TYPE = 0x7D
	F64_NUM_TYPE = 0x7C
)

type Module struct {
	Magic   uint32 // magic number >`\0asm`
	Version uint32 // version
}
