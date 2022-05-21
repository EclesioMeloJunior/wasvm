package parser

import "fmt"

func BinaryFormat(filepath string) error {
	bp, err := NewBinaryParser(filepath)
	if err != nil {
		return err
	}

	// starting parsing the `wasm header` values
	if err = bp.ParseMagicNumber(); err != nil {
		return fmt.Errorf("cannot parse magic number: %w", err)
	}

	if err = bp.ParseVersion(); err != nil {
		return fmt.Errorf("cannot parse version: %w", err)
	}

	if err = bp.ParseSection(); err != nil {
		return fmt.Errorf("cannot parse section: %w", err)
	}

	return nil
}
