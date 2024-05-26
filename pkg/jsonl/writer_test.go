package jsonl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter_Pos_WriteJson(t *testing.T) {
	testStructs := []testStruct{
		{Name: "Gilbert", Wins: [][]string{{"straight", "7♣"}, {"one pair", "10♥"}}},
		{Name: "Alexa", Wins: [][]string{{"two pair", "4♠"}, {"two pair", "9♠"}}},
		{Name: "May", Wins: [][]string{}},
		{Name: "Deloise", Wins: [][]string{{"three of a kind", "5♣"}}},
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(testString)))
	wr := NewWriter(buf)

	for _, st := range testStructs {
		assert.NoError(t, wr.WriteJson(st))
	}

	assert.NoError(t, wr.Flush())
	assert.Equal(t, testString, buf.String())
}

func TestWriter_Neg_WriteJson(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, len(testString)))
	wr := NewWriter(buf)

	assert.Error(t, wr.WriteJson(func() {}))
}
