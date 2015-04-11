package gotorrent

import "testing"

type TestList struct {
	Name string
	Age  int
}

func ValidateMarshal(input interface{}, expected string, t *testing.T) {
	actual, err := Marshal(input)
	if err != nil {
		t.Errorf("Error marshalling %v (expected %v, received %v): %v",
			input, expected, actual, err)
	}
	if expected != actual {
		t.Errorf("Expected %v, got %v on input %v", expected, actual, input)
	}
}

func TestMarshal(t *testing.T) {
	marshalData := map[interface{}]string{
		10:                    "i10e",
		"spam":                "4:spam",
		TestList{"alice", 30}: "l5:alicei30ee",
	}
	for input, expected := range marshalData {
		ValidateMarshal(input, expected, t)
	}
	ValidateMarshal(
		map[TestList]int{TestList{"bob", 25}: 30, TestList{"alice", 30}: 35},
		"dl3:bobi25eei30el5:alicei30eei35ee", t)
}
