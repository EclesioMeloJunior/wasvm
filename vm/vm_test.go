package vm_test

import (
	"fmt"
	"testing"

	"github.com/EclesioMeloJunior/wasvm/parser"
	"github.com/EclesioMeloJunior/wasvm/vm"
	"github.com/stretchr/testify/assert"
)

const simpleWasm = "../tests/simple.wasm"

func TestSimpleWasm_ExportedFunction_Execution(t *testing.T) {
	wasm, err := parser.NewBinaryParser(simpleWasm)
	assert.NoError(t, err)

	rt, err := vm.NewRuntime(wasm)
	assert.NoError(t, err)

	// TODO: check why the exported function is not exported
	fmt.Println(rt.Exported)
}
