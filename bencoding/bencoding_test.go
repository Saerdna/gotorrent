package bencoding

import (
	"bytes"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	testData := map[int]string{
		10:    "i10e",
		0:     "i0e",
		12345: "i12345e",
		-100:  "i-100e",
	}
	for input, expected := range testData {
		var buffer bytes.Buffer
		EncodeInt(input, buffer)
		actual := buffer.String()
		if expected != actual {
			t.Errorf("Expected: %v Actual: %v\n", expected, actual)
		}
	}
}
