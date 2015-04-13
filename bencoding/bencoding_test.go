package bencoding

import (
	"bufio"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
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
	tokenReader := TokenReader{bufio.NewReader(strings.NewReader(input))}
	for i, expectedToken := range expectedTokens {
		expectedValue := expectedValues[i]
		actualToken, actualValue := tokenReader.NextToken()
		if actualToken != expectedToken {
			t.Errorf("Expected %v, Actual %v", expectedToken, actualToken)
		}
		if actualValue != expectedValue {
			t.Errorf("Expected %v, Actual %v", expectedValue, actualValue)
		}
	}
}

func TestParser(t *testing.T) {
	input := "d4:name5:alice4:kidsl3:bob5:carole3:agei30ee"
	expected := Dict{map[Node]Node{
		String{"name"}: String{"alice"},
		String{"kids"}: List{[]Node{
			String{"bob"},
			String{"carol"},
		}},
		String{"age"}: Int{30},
	}}
	actual, err := ParseString(input)
	if err != nil {
		t.Errorf("Error parsing %v: %v", input, err)
	}
	if !Equals(actual, expected) {
		t.Errorf("Expected %v, Actual %v", expected, actual)
	}
}
