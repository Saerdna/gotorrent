package gotorrent

import (
	"fmt"
	"strconv"
	"strings"
)

type TypeBencodedData interface {
	isBencoded()
}

type TypeBencodedString string
type TypeBencodedInteger string
type TypeBencodedList string

func (TypeBencodedString) isBencoded()  {}
func (TypeBencodedInteger) isBencoded() {}
func (TypeBencodedList) isBencoded()    {}

type TypeBdecodedData interface {
	isBdecoded()
}

type TypeBdecodedString string
type TypeBdecodedInt int
type TypeBdecodedList []TypeBdecodedData

func (TypeBdecodedString) isBdecoded() {}
func (TypeBdecodedInt) isBdecoded() {}
func (TypeBdecodedList) isBdecoded() {}

func BencodeString(s TypeBdecodedString) TypeBencodedString {
	return TypeBencodedString(strconv.Itoa(len(string(s)) + ":" + string(s))
}

func BdecodeString(b TypeBencodedString) (TypeBdecodedString, error) {
	substrings := strings.SplitN(string(b), ":", 2)
	if len(substrings) < 2 {
		return TypeBdecodedString(""), fmt.Errorf("Separator not found %v", b)
	}
	string_length, err := strconv.Atoi(substrings[0])
	if err != nil || len(substrings[1]) != string_length {
		return TypeBdecodedString(""), fmt.Errorf("Invalid length encoding %v", b)
	}
	return TypeBdecodedString(substrings[1]), nil
}

func BencodeInteger(i TypeBdecodedInt) TypeBencodedInteger {
	return TypeBencodedInteger("i" + strconv.Itoa(int(i)) + "e")
}

func BdecodeInteger(b TypeBencodedInteger) (TypeBdecodedInt, error) {
	s := string(b)
	if !strings.HasPrefix(s, "i") || !strings.HasSuffix(s, "e") {
		return TypeBdecodedInt(0), fmt.Errorf("Invalid bencoded integer %v", b)
	}
	result, err := strconv.Atoi(s[1 : len(s)-1])
	if err != nil {
		return TypeBdecodedInt(0), fmt.Errorf("Unable to decode integer %v", b)
	}
	return TypeBdecodedInt(result), nil
}

func BencodeList(l TypeBdecodedList) TypeBencodedList {
	data := make([]string, len(l))
	for i := 0; i < len(l); i++ {
		switch t := l[i].(type) {
		case TypeBencodedString:
			data[i] = string(TypeBencodedString(t))
		case TypeBencodedInteger:
			data[i] = string(TypeBencodedInteger(t))
		case TypeBencodedList:
			data[i] = string(TypeBencodedList(t))
		}
	}
	return TypeBencodedList("l" + strings.Join(data, "") + "e")
}

func BdecodeList(l TypeBencodedList) TypeBdecodedList, error {
	s := string(l)
	if !strings.HasPrefix(s, "l") || !strings.HasSuffix(s, "e") {
		return nil, fmt.Errorf("Invalid bencoded list %v", l)
	}
}
