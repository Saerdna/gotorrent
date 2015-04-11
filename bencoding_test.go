package gotorrent

import "testing"

func TestBencodeString(t *testing.T) {
	expected := TypeBencodedString("4:test")
	actual := BencodeString("test")
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestBdecodeString(t *testing.T) {
	expected := "test"
	actual, err := BdecodeString(TypeBencodedString("4:test"))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestBdecodeMalformedStrings(t *testing.T) {
	_, err := BdecodeString("test")
	if err == nil {
		t.Errorf("Failed to detect error in %v", "test")
	}
	_, err = BdecodeString("a:test")
	if err == nil {
		t.Errorf("Failed to detect error in %v", "test")
	}
}

func TestBencodeInteger(t *testing.T) {
	expected := TypeBencodedInteger("i10e")
	actual := BencodeInteger(10)
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestBdecodeInteger(t *testing.T) {
	expected := 10
	actual, err := BdecodeInteger(TypeBencodedInteger("i10e"))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestBencodeList(t *testing.T) {
	expected := TypeBencodedList("l4:spam4:eggsi10ee")
	data := []TypeBencodedData{BencodeString("spam"), BencodeString("eggs"), BencodeInteger(10)}
	actual := BencodeList(data)
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestBencodeSubList(t *testing.T) {
	s := "l4:spam4:eggsi10ee"
	expected := TypeBencodedList("l" + s + s + "e")
	data := []TypeBencodedData{BencodeString("spam"), BencodeString("eggs"), BencodeInteger(10)}
	data2 := []TypeBencodedData{BencodeList(data), BencodeList(data)}
	actual := BencodeList(data2)
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
