package bencoding

import (
	"bytes"
	"strconv"
)

func EncodeInt(i int, writer bytes.Buffer) {
	writer.WriteString("i")
	writer.WriteString(strconv.Itoa(i))
	writer.WriteString("e")
}

func EncodeString(s string, writer bytes.Buffer) {
	writer.WriteString(strconv.Itoa(len(s)))
	writer.WriteString(":")
	writer.WriteString(s)
}

func EncodeList(l []string, writer bytes.Buffer) {
	writer.WriteString("l")
	for _, element := range l {
		writer.WriteString(element)
	}
	writer.WriteString("e")
}

func EncodeDict(m map[string]string, writer bytes.Buffer) {
	writer.WriteString("d")
	for key, value := range m {
		writer.WriteString(key)
		writer.WriteString(value)
	}
	writer.WriteString("e")
}
