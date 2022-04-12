package gowasm_test

import (
	"fmt"
	"testing"

	"github.com/EclesioMeloJunior/gowasm"
	"github.com/stretchr/testify/assert"
)

const (
	SIMPLE_WASM    = "./.wasms/simple.wasm"
	FIBONACCI_WASM = "./.wasms/fibonacci.wasm"
)

func TestDecoderLoad(t *testing.T) {
	contents, err := gowasm.Load(SIMPLE_WASM, gowasm.BinaryFormat)
	assert.NoError(t, err)

	for _, c := range contents {
		fmt.Printf("0x%x\n", c)
	}
}
