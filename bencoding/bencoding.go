package bencoding

import (
	"bufio"
	"strconv"
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

func NextToken(b *bufio.Reader) (token Token, value string) {
	r, _, err := b.ReadRune()
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
		i, err := b.ReadString('e')
		if err != nil {
			return EOF, ""
		}
		return INT, i[0 : len(i)-1]
	case '0' <= r && r <= '9':
		b.UnreadRune()
		length_string, err := b.ReadString(':')
		if err != nil {
			return EOF, ""
		}

		length, err := strconv.Atoi(length_string[:len(length_string)-1])
		if err != nil {
			return ILLEGAL, ""
		}

		buffer := make([]byte, length)
		n, err := b.Read(buffer)
		if err != nil || n != length {
			return EOF, ""
		}
		return STRING, string(buffer)
	default:
		return ILLEGAL, ""
	}
}

// type Node interface {
// 	Decode(*bufio.Reader) error
// }

// type Int struct {
// 	Value int
// }

// type String struct {
// 	Value string
// }

// type List struct {
// 	Value []Node
// }

// type End struct{}

// func DecodeUnknown(b *bufio.Reader) (Node, error) {
// 	start, err := b.Peek(1)
// 	if err != nil {
// 		return err
// 	}
// 	for {
// 		switch string(start) {
// 		case "i":
// 			var value Int
// 			err = value.Decode(r)
// 			if err != nil {
// 				return value, err
// 			}
// 			return value, nil
// 		case "l":
// 			var value List
// 			err = value.Decode(r)
// 			if err != nil {
// 				return value, err
// 			}
// 			return value, nil
// 		case "d":
// 			var value Dict
// 			err = value.Decode(r)
// 			if err != nil {
// 				return value, err
// 			}
// 			return value, nil
// 		case "e":
// 			return End{}, nil
// 		default:
// 			var value String
// 			err = value.Decode(r)
// 			if err != nil {
// 				return value, err
// 			}
// 			return value, nil
// 		}
// 	}
// }

// func (i *Int) Decode(b *bufio.Reader) error {
// 	start, err := b.ReadByte()
// 	if err != nil {
// 		return err
// 	}
// 	if start != 'i' {
// 		return fmt.Errorf("Int didn't start with i.")
// 	}

// 	s, err := b.ReadString('e')
// 	if err != nil {
// 		return err
// 	}

// 	i.Value, err = strconv.Atoi(s[:len(s)-1])
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *String) Decode(b *bufio.Reader) error {
// 	length_string, err := b.ReadString(':')
// 	if err != nil {
// 		return err
// 	}

// 	length, err := strconv.Atoi(length_string[:len(length_string)-1])
// 	if err != nil {
// 		return err
// 	}

// 	buffer := make([]byte, length)
// 	n, err := b.Read(buffer)
// 	if err != nil {
// 		return err
// 	}
// 	if n != length {
// 		return fmt.Errorf("String had length %v, but buffer only had %v bytes: %v",
// 			length, len(s.Value), s.Value)
// 	}

// 	s.Value = string(buffer)
// 	return nil
// }

// func (l *List) Decode(b *bufio.Reader) error {
// 	start := make([]byte, 1)
// 	n, err := b.Read(start)
// 	if err != nil {
// 		return err
// 	}
// 	if n < len(start) || string(start) != "l" {
// 		return fmt.Errorf("Invalid list.")
// 	}

// 	start, err = b.Peek(1)
// 	if err != nil {
// 		return err
// 	}
// 	for {
// 		switch string(start) {
// 		case "i":
// 			var value Int
// 			err = value.Decode(r)
// 			if err != nil {
// 				return err
// 			}
// 			l.Value = append(l.Value, value)
// 		case "l":
// 			var value List
// 			err = value.Decode(r)
// 			if err != nil {
// 				return err
// 			}
// 			l.Value = append(l.Value, value)
// 		case "d":
// 			var value Dict
// 			err = value.Decode(r)
// 			if err != nil {
// 				return err
// 			}
// 			l.Value = append(l.Value, value)
// 		default:
// 			var value String
// 			err = value.Decode(r)
// 			if err != nil {
// 				return err
// 			}
// 			l.Value = append(l.Value, value)
// 		}
// 	}
// }

// func (e *End) Decode(b *bufio.Reader) error {

// }
