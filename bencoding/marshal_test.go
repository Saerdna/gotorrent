package bencoding

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
		return
	}
	if expected != actual {
		t.Errorf("Expected %v, got %v on input %v", expected, actual, input)
		return
	}
}

func TestMarshal(t *testing.T) {
	marshalData := map[interface{}]string{
		10:                    "i10e",
		"spam":                "4:spam",
		TestList{"alice", 30}: "d4:name5:alice3:agei30ee",
	}
	for input, expected := range marshalData {
		ValidateMarshal(input, expected, t)
	}
	ValidateMarshal(
		map[TestList]int{TestList{"alice", 30}: 35, TestList{"bob", 25}: 30},
		"dd4:name3:bob3:agei25eei30ed4:name5:alice3:agei30eei35ee", t)
	ValidateMarshal([]int{10, 20, 30}, "li10ei20ei30ee", t)
}
