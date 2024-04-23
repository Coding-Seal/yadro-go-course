package jsonl

import (
	"bytes"
	"testing"
)

func TestWriter_WriteJson(t *testing.T) {
	testStructs := []testStruct{
		testStruct{Name: "Gilbert", Wins: [][]string{{"straight", "7♣"}, {"one pair", "10♥"}}},
		testStruct{Name: "Alexa", Wins: [][]string{{"two pair", "4♠"}, {"two pair", "9♠"}}},
		testStruct{Name: "May", Wins: [][]string{}},
		testStruct{Name: "Deloise", Wins: [][]string{{"three of a kind", "5♣"}}},
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(testString)))
	wr := NewWriter(buf)

	for _, st := range testStructs {
		err := wr.WriteJson(st)
		if err != nil {
			t.Error(err)
		}
	}

	if err := wr.Flush(); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(buf.Bytes(), []byte(testString)) {
		t.Errorf("%s != %s", buf.Bytes(), []byte(testString))
	}
}
