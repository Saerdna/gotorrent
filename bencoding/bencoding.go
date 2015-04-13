package bencoding

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Token int

const (
	EOF Token = iota
	LIST_START
	DICT_START
	END
	INT
	STRING
	ILLEGAL
)

type TokenReader struct {
	b *bufio.Reader
}

func (t *TokenReader) NextToken() (token Token, value string) {
	r, _, err := t.b.ReadRune()
	if err != nil {
		return EOF, ""
	}
	switch {
	case r == 'l':
		return LIST_START, ""
	case r == 'd':
		return DICT_START, ""
	case r == 'e':
		return END, ""
	case r == 'i':
		i, err := t.b.ReadString('e')
		if err != nil {
			return EOF, ""
		}
		return INT, i[0 : len(i)-1]
	case '0' <= r && r <= '9':
		t.b.UnreadRune()
		length_string, err := t.b.ReadString(':')
		if err != nil {
			return EOF, ""
		}

		length, err := strconv.Atoi(length_string[:len(length_string)-1])
		if err != nil {
			return ILLEGAL, ""
		}

		buffer := make([]byte, length)
		n, err := t.b.Read(buffer)
		if err != nil || n != length {
			return EOF, ""
		}
		return STRING, string(buffer)
	default:
		return ILLEGAL, ""
	}
}

type Node interface {
	isNode()
}

type Int struct {
	Int int
}
type String struct {
	String string
}
type Dict struct {
	Dict map[Node]Node
}
type List struct {
	List []Node
}

func (Int) isNode()    {}
func (String) isNode() {}
func (Dict) isNode()   {}
func (List) isNode()   {}

func ParseString(s string) (Node, error) {
	tokenReader := TokenReader{bufio.NewReader(strings.NewReader(s))}
	return Parse(&tokenReader)
}
func Parse(t *TokenReader) (Node, error) {
	token, value := t.NextToken()
	switch token {
	case EOF:
		return nil, nil
	case LIST_START:
		l := []Node{}
		for {
			value, err := Parse(t)
			if err != nil {
				return nil, err
			}
			if value == nil {
				break
			}
			l = append(l, value)
		}
		return List{l}, nil
	case DICT_START:
		m := map[Node]Node{}
		for {
			key, err := Parse(t)
			if err != nil {
				return nil, err
			}
			if key == nil {
				break
			}
			value, err := Parse(t)
			m[key] = value
		}
		return Dict{m}, nil
	case END:
		return nil, nil
	case INT:
		i, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return Int{i}, nil
	case STRING:
		return String{value}, nil
	case ILLEGAL:
		return nil, fmt.Errorf("Illegal character in stream")
	default:
		return nil, fmt.Errorf("Unknown tokent: %v", token)
	}
}

func Equals(n1, n2 Node) bool {
	switch t1 := n1.(type) {
	case Int:
		t2, ok := n2.(Int)
		if !ok || t1.Int != t2.Int {
			return false
		}
		return true
	case String:
		t2, ok := n2.(String)
		if !ok || t1.String != t2.String {
			return false
		}
		return true
	case List:
		t2, ok := n2.(List)
		if !ok || len(t1.List) != len(t2.List) {
			return false
		}
		for i, value1 := range t1.List {
			value2 := t2.List[i]
			if !Equals(value1, value2) {
				return false
			}
		}
		return true
	case Dict:
		t2, ok := n2.(Dict)
		if !ok {
			return false
		}
		if len(t1.Dict) != len(t2.Dict) {
			return false
		}
		for key, value1 := range t1.Dict {
			value2, ok := t2.Dict[key]
			if !ok || !Equals(value1, value2) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
