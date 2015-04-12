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
		map[TestList]int{TestList{"alice", 30}: 35, TestList{"bob", 25}: 30},
		"dl3:bobi25eei30el5:alicei30eei35ee", t)
	ValidateMarshal([]int{10, 20, 30}, "li10ei20ei30ee", t)
}

func ValidateUnmarshal(input string, expected, actual interface{}, err error, t *testing.T) {
	if err != nil {
		t.Errorf("Error unmarshalling %v (expected %v, received %v): %v",
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
	var actual []string
	input := "l4:spam4:eggs5:bacone"
	expected := []string{"spam", "eggs", "bacon"}
	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("Error unmarshalling %v (expected %v, received %v): %v",
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
	var actual []int
	input := "li10ei20ei30ee"
	expected := []int{10, 20, 30}
	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("Error unmarshalling %v (expected %v, received %v): %v",
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

func TestUnmarshalArrayOfArrays(t *testing.T) {
	var actual [][]int
	input := "l" + "li10ei20ei30ee" + "li10ei20ee" + "e"
	expected := [][]int{{10, 20, 30}, {10, 20}}
	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("Error unmarshalling %v (expected %v, received %v): %v",
			input, expected, actual, err)
		return
	}
	if len(expected) != len(actual) {
		t.Errorf("Different lengths.  Expected: %v Actual: %v Input: %v", expected, actual, input)
		return
	}
	for i := 0; i < len(expected); i++ {
		if len(expected[i]) != len(actual[i]) {
			t.Errorf("Different lengths.  Expected: %v Actual: %v Input: %v", expected, actual, input)
			return
		}
		for j := 0; j < len(expected[i]); j++ {
			if expected[i][j] != actual[i][j] {
				t.Errorf("Expected %v, got %v on input %v", expected[i][j], actual[i][j], input)
			}
		}
	}
}

func TestUnmarshalStruct(t *testing.T) {
	actual := TestList{}
	input := "l5:alicei30ee"
	expected := TestList{"alice", 30}
	err := Unmarshal(input, &actual)
	ValidateUnmarshal(input, expected, actual, err, t)
}

type TestListWithPrivate struct {
	Name  string
	dummy string
	Age   int
}

func TestUnmarshalStructSkipsPrivate(t *testing.T) {
	actual := TestListWithPrivate{}
	input := "l5:alicei30ee"
	expected := TestListWithPrivate{Name: "alice", Age: 30}
	err := Unmarshal(input, &actual)
	ValidateUnmarshal(input, expected, actual, err, t)
}

type TestListWithArray struct {
	Name string
	Kids []string
	Age  int
}

func TestUnmarshalStructWithSlice(t *testing.T) {
	actual := TestListWithArray{}
	input := "l5:alicel3:bob5:carolei30ee"
	expected := TestListWithArray{"alice", []string{"bob", "carol"}, 30}
	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("Error unmarshalling %v (expected %v, received %v): %v",
			input, expected, actual, err)
		return
	}
	if expected.Name != actual.Name || expected.Age != actual.Age {
		t.Errorf("Expected %v, got %v on input %v", expected, actual, input)
		return
	}
	if len(actual.Kids) != len(expected.Kids) {
		t.Errorf("Different lengths.  Expected: %v Actual: %v Input: %v",
			expected.Kids, actual.Kids, input)
		return
	}
	for i := 0; i < len(actual.Kids); i++ {
		if expected.Kids[i] != actual.Kids[i] {
			t.Errorf("Expected %v, got %v on input %v", expected.Kids[i], actual.Kids[i], input)
			return
		}
	}
}

func TestUnmarshalMap(t *testing.T) {
	actual := map[string]int{}
	input := "d5:alicei30e3:bobi25ee"
	expected := map[string]int{
		"alice": 30,
		"bob":   25,
	}
	err := Unmarshal(input, &actual)
	if err != nil {
		t.Errorf("Error unmarshalling %v (expected %v, received %v): %v",
			input, expected, actual, err)
		return
	}
	if len(expected) != len(actual) {
		t.Errorf("Different lengths.  Expected: %v Actual: %v Input: %v", expected, actual, input)
	}
	for k, v := range expected {
		val, present := actual[k]
		if !present {
			t.Errorf("Key not found in actual: %v", k)
		}
		if val != v {
			t.Errorf("Different values for key %v.  Expected: %v Actual: %v", k, v, val)
		}
	}
}
