package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIFOpCodeIntruction(t *testing.T) {
	tests := map[string]struct {
		instructions []byte
		wantErr      string
		expected     []any
		results      []any
	}{
		"does not have end if": {
			instructions: []byte{
				0x04, 0x7F, 0x41, 0x01,
			},
			wantErr: "failed to find if end",
		},
		"only if + end": {
			instructions: []byte{
				0x41, 0x01, // put 02 in the stack
				0x41, 0x02, // put 01 in the stack
				0x48,                         // 02 > 01 (true)
				0x04, 0x7F, 0x41, 0x01, 0x0B, // if condition + if end
				0x0B, // function end
			},
			expected: []any{int32(1)}, // we spect the number 1 only
			results:  []any{int32(0)}, //define the result type
		},
		"if + else + end, testing the if branch": {
			instructions: []byte{
				0x41, 0x01, // put 01 in the stack
				0x41, 0x02, // put 02 in the stack
				0x48,                   // 02 > 01 (true)
				0x04, 0x7F, 0x41, 0x01, // if condition
				0x05, 0x041, 0x02, 0x0B, // else condition + if end
				0x0B, // function end
			},
			expected: []any{int32(1)}, // we spect the number 1 only
			results:  []any{int32(0)}, //define the result type
		},
		"if + else + end, testing the else branch": {
			instructions: []byte{
				0x41, 0x02, // put 02 in the stack
				0x41, 0x01, // put 01 in the stack
				0x48,                   // 01 > 02 (false)
				0x04, 0x7F, 0x41, 0x01, // if condition
				0x05, 0x041, 0x02, 0x0B, // else condition + if end
				0x0B, // function end
			},
			expected: []any{int32(2)}, // we spect the number 1 only
			results:  []any{int32(0)}, //define the result type
		},
		"nested ifs": {
			instructions: []byte{
				0x41, 0x01, // put 01 in the stack
				0x41, 0x02, // put 02 in the stack
				0x48,       // 02 > 01 (true)
				0x04, 0x7F, // if condition
				0x41, 0x01, // put 01 in the stack
				0x41, 0x02, // put 02 in the stack
				0x48,
				0x04, 0x7F, // another if condition
				0x41, 0x03, // put 03 in the stack
				0x41, 0x04, // put 04 in the stack
				0x6A, // sum them up and return
				0x0B, // end nested if
				0x0B, // end if
			},
			expected: []any{int32(7)}, // we spect the number 1 only
			results:  []any{int32(0)}, //define the result type
		},
	}

	for tname, tt := range tests {
		tt := tt
		t.Run(tname, func(t *testing.T) {
			cf := &callFrame{
				pc:           0,
				stack:        make([]StackValue, 0, 1024),
				instructions: tt.instructions,
				params:       []any{},
				results:      []any{int32(0)},
			}

			res, err := cf.Call()
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}
