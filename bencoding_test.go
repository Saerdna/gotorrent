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
		TestList{"alice", 30}: "l5:alicei30ee",
	}
	for input, expected := range marshalData {
		ValidateMarshal(input, expected, t)
	}
	ValidateMarshal(
		map[TestList]int{TestList{"bob", 25}: 30, TestList{"alice", 30}: 35},
		"dl3:bobi25eei30el5:alicei30eei35ee", t)
	ValidateMarshal([]int{10, 20, 30}, "li10ei20ei30ee", t)
}

func ValidateUnmarshal(input string, expected, actual interface{}, err error, t *testing.T) {
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

func TestUnmarshalInt(t *testing.T) {
	var actual int
	input := "i10e"
	expected := 10
	err := Unmarshal(input, &actual)
	ValidateUnmarshal(input, expected, actual, err, t)
}

func TestUnmarshalString(t *testing.T) {
	var actual string
	input := "4:spam"
	expected := "spam"
	err := Unmarshal(input, &actual)
	ValidateUnmarshal(input, expected, actual, err, t)
}

func TestUnmarshalStringArray(t *testing.T) {
	actual := make([]string, 3)
	input := "l4:spam4:eggs5:bacone"
	expected := []string{"spam", "eggs", "bacon"}
	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("Error marshalling %v (expected %v, received %v): %v",
			input, expected, actual, err)
		return
	}
	if len(expected) != len(actual) {
		t.Errorf("Different lengths.  Expected: %v Actual: %v Input: %v", expected, actual, input)
		return
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			t.Errorf("Expected %v, got %v on input %v", expected[i], actual[i], input)
		}
	}
}

func TestUnmarshalIntegerArray(t *testing.T) {
	actual := make([]int, 3)
	input := "li10ei20ei30ee"
	expected := []int{10, 20, 30}
	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("Error marshalling %v (expected %v, received %v): %v",
			input, expected, actual, err)
		return
	}
	if len(expected) != len(actual) {
		t.Errorf("Different lengths.  Expected: %v Actual: %v Input: %v", expected, actual, input)
		return
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			t.Errorf("Expected %v, got %v on input %v", expected[i], actual[i], input)
		}
	}
}
