package scroll

import (
	"bytes"
	"testing"
)

func TestArraysToBytes(t *testing.T) {
	// wrongly sized arrays
	if ArraysToBytes([][]uint8{{0, 0}, {1, 1}}) != nil {
		t.Errorf("Argument for ArraysToBytes() should be in size: %dx%d.", NumRows, NumCols)
	}

	// hand-written arrays
	arg := [][]uint8{
		{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
	}
	converted := ArraysToBytes(arg)
	expected := []byte{31, 0, 31, 0, 31, 0, 31, 0, 31, 0, 31}
	if !bytes.Equal(converted, expected) {
		t.Errorf("ArraysToBytes(%+v) returned %+v, expected: %+v.", arg, converted, expected)
	}
}
