package scroll

import (
	"bytes"
	"testing"
)

func TestTrim(t *testing.T) {
	var b, trimmed, expected []byte

	// all blank => two spaces
	b = []byte{0, 0, 0, 0, 0, 0}
	trimmed = trim(b)
	expected = []byte{0, 0}
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("trim(%+v) returned: %+v, expected: %+v.", b, trimmed, expected)
	}

	// trim spaces on left side
	b = []byte{0, 0, 1, 1, 0, 0, 1, 1}
	trimmed = trim(b)
	expected = []byte{1, 1, 0, 0, 1, 1}
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("trim(%+v) returned: %+v, expected: %+v.", b, trimmed, expected)
	}

	// trim spaces on right side
	b = []byte{1, 1, 0, 0, 1, 1, 0, 0}
	trimmed = trim(b)
	expected = []byte{1, 1, 0, 0, 1, 1}
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("trim(%+v) returned: %+v, expected: %+v.", b, trimmed, expected)
	}

	// trim spaces on both sides
	b = []byte{0, 0, 1, 1, 0, 0, 1, 1, 0, 0}
	trimmed = trim(b)
	expected = []byte{1, 1, 0, 0, 1, 1}
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("trim(%+v) returned: %+v, expected: %+v.", b, trimmed, expected)
	}
}

func TestBytesForCharacter(t *testing.T) {
	var c rune
	var converted, expected []byte

	// unconvertible character => spaces
	c = '헐'
	converted = BytesForCharacter(c)
	expected = []byte{0, 0}
	if !bytes.Equal(converted, expected) {
		t.Errorf("BytesForCharacter(%+v) returned: %+v, expected: %+v. (unconvertible character should be converted to spaces)", c, converted, expected)
	}

	// convertible character
	c = 'a'
	converted = BytesForCharacter(c)
	expected = []byte{8, 20, 28}
	if !bytes.Equal(converted, expected) {
		t.Errorf("BytesForCharacter(%+v) returned: %+v, expected: %+v.", c, converted, expected)
	}
}

func TestBytesForString(t *testing.T) {
	var s string
	var converted, expected []byte

	// blank string
	s = "    "
	converted = BytesForString(s)
	expected = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(converted, expected) {
		t.Errorf("BytesForString(%+v) returned: %+v, expected: %+v.", s, converted, expected)
	}

	// ordinary string
	s = "Go go go."
	converted = BytesForString(s)
	expected = []byte{14, 17, 29, 0, 8, 20, 8, 0, 0, 0, 0, 2, 21, 15, 0, 8, 20, 8, 0, 0, 0, 0, 2, 21, 15, 0, 8, 20, 8, 0, 16, 0}
	if !bytes.Equal(converted, expected) {
		t.Errorf("BytesForString(%+v) returned: %+v, expected: %+v.", s, converted, expected)
	}
}
