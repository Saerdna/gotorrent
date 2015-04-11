package gotorrent

import "testing"

func TestMarshal(t *testing.T) {
	marshalData := map[interface{}]string{
		10:     "i10e",
		"spam": "4:spam",
	}
	for input, expected := range marshalData {
		actual, err := Marshal(input)
		if err != nil {
			t.Errorf("Error marshalling %v (expected %v, received %v): %v",
				input, expected, actual, err)
		}
		if expected != actual {
			t.Errorf("Expected %v, got %v on input %v", expected, actual, input)
		}
	}
}
