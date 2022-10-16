package vm_test

import (
	"testing"

	"github.com/EclesioMeloJunior/wasvm/parser"
	"github.com/EclesioMeloJunior/wasvm/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const simpleWasm = "../resources/simple.wasm"
const operationsWasm = "../resources/operations.wasm"
const factorialWasm = "../resources/factorial.wasm"

func TestSimpleWasm_ExportedFunction_Execution(t *testing.T) {
	binaryWASM, err := parser.BinaryFormat(simpleWasm)
	assert.NoError(t, err)

	rt, err := vm.NewRuntime(binaryWASM)
	assert.NoError(t, err)

	assert.Len(t, rt.Exported, 1)

	const exportedFun = "helloWorld"
	callFrame, ok := rt.Exported[exportedFun]
	assert.True(t, ok)
	assert.NotNil(t, callFrame)

	results, err := callFrame.Call()
	assert.NoError(t, err)

	require.Len(t, results, 1)
	result0 := results[0]
	value, ok := result0.(int32)

	require.True(t, ok)
	require.Equal(t, int32(42), value)
}

func TestOperationsWasm(t *testing.T) {
	binaryWASM, err := parser.BinaryFormat(operationsWasm)
	assert.NoError(t, err)

	rt, err := vm.NewRuntime(binaryWASM)
	assert.NoError(t, err)
	assert.Len(t, rt.Exported, 3)

	tests := []struct {
		function string
		lhs      int32
		rhs      int32
		expected int32
	}{
		{
			function: "sum",
			lhs:      10,
			rhs:      10,
			expected: 20,
		},
		{
			function: "sub",
			lhs:      0,
			rhs:      1,
			expected: -1,
		},
		{
			function: "mul",
			lhs:      9,
			rhs:      8,
			expected: 72,
		},

		//TODO: implement div operation
	}

	for _, tt := range tests {
		function, ok := rt.Exported[tt.function]
		assert.True(t, ok)

		results, err := function.Call(tt.lhs, tt.rhs)
		assert.NoError(t, err)
		require.Len(t, results, 1)

		result, ok := results[0].(int32)
		require.True(t, ok)
		assert.Equal(t, tt.expected, result)
	}
}

func TestFactorialWasm(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		param, expected int32
	}{
		"factorial_3": {
			param: 3, expected: 6,
		},
		"factorial_5": {
			param: 5, expected: 120,
		},
		"factorial_10": {
			param: 10, expected: 3628800,
		},
		"factorial_100": {
			param: 12, expected: 479001600,
		},
	}

	for tname, tt := range tests {
		tt := tt
		t.Run(tname, func(t *testing.T) {
			t.Parallel()

			binaryWASM, err := parser.BinaryFormat(factorialWasm)
			assert.NoError(t, err)

			rt, err := vm.NewRuntime(binaryWASM)
			assert.NoError(t, err)
			assert.Len(t, rt.Exported, 1)

			fac, ok := rt.Exported["fac"]
			assert.True(t, ok)

			results, err := fac.Call(tt.param)
			assert.NoError(t, err)
			require.Len(t, results, 1)

			assert.Equal(t, results[0], tt.expected)
		})
	}
}
