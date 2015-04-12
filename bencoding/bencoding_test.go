package bencoding

import (
	"bufio"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	input := "d4:name5:alice4:kidsl3:bob5:carole3:agei30ee"
	expectedTokens := []Token{
		DICT_START,
		STRING,
		STRING,
		STRING,
		LIST_START,
		STRING,
		STRING,
		END,
		STRING,
		INT,
		END,
	}
	expectedValues := []string{
		"",
		"name",
		"alice",
		"kids",
		"",
		"bob",
		"carol",
		"",
		"age",
		"30",
		"",
	}
	reader := bufio.NewReader(strings.NewReader(input))
	for i, expectedToken := range expectedTokens {
		expectedValue := expectedValues[i]
		actualToken, actualValue := NextToken(reader)
		if actualToken != expectedToken {
			t.Errorf("Expected %v, Actual %v", expectedToken, actualToken)
		}
		if actualValue != expectedValue {
			t.Errorf("Expected %v, Actual %v", expectedValue, actualValue)
		}
	}
}

// func TestDecodeInt(t *testing.T) {
// 	actual := Int{}
// 	input := "i10e"
// 	expected := 10
// 	err := actual.Decode(strings.NewReader(input))
// 	if err != nil {
// 		t.Errorf("Error decoding %v %v", input, err)
// 	}
// 	if actual.Value != expected {
// 		t.Errorf("Results don't match when decoding %v.  Expected %v Actual %v",
// 			input, expected, actual)
// 	}
// }

// func TestDecodeString(t *testing.T) {
// 	actual := String{}
// 	input := "4:spam"
// 	expected := "spam"
// 	err := actual.Decode(strings.NewReader(input))
// 	if err != nil {
// 		t.Errorf("Error decoding %v %v", input, err)
// 	}
// 	if actual.Value != expected {
// 		t.Errorf("Results don't match when decoding %v.  Expected %v Actual %v",
// 			input, expected, actual)
// 	}
// }
