package binary_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/EclesioMeloJunior/gowasm/binary"
	"github.com/EclesioMeloJunior/gowasm/collections"
	"github.com/stretchr/testify/assert"
)

const simpleWasm = "../tests/simple.wasm"

func TestBinaryRead(t *testing.T) {
	var iter collections.Iterator[byte]
	iter, err := binary.NewWasm(simpleWasm)
	assert.NoError(t, err)

	for {
		b, err := iter.Next()
		if err != nil {
			assert.ErrorIs(t, err, io.EOF)
			break
		}

		fmt.Printf("0x%x\n", b)
	}

}
