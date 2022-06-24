package vm_test

import (
	"testing"

	"github.com/EclesioMeloJunior/wasvm/parser"
	"github.com/EclesioMeloJunior/wasvm/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const simpleWasm = "../tests/simple.wasm"

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
