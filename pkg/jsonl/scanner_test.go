package jsonl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testString string = "{\"name\":\"Gilbert\",\"wins\":[[\"straight\",\"7♣\"],[\"one pair\",\"10♥\"]]}\n{\"name\":\"Alexa\",\"wins\":[[\"two pair\",\"4♠\"],[\"two pair\",\"9♠\"]]}\n{\"name\":\"May\",\"wins\":[]}\n{\"name\":\"Deloise\",\"wins\":[[\"three of a kind\",\"5♣\"]]}\n"

type testStruct struct {
	Name string     `json:"name"`
	Wins [][]string `json:"wins"`
}

func TestScanner(t *testing.T) {
	expectedRes := []testStruct{
		{Name: "Gilbert", Wins: [][]string{{"straight", "7♣"}, {"one pair", "10♥"}}},
		{Name: "Alexa", Wins: [][]string{{"two pair", "4♠"}, {"two pair", "9♠"}}},
		{Name: "May", Wins: [][]string{}},
		{Name: "Deloise", Wins: [][]string{{"three of a kind", "5♣"}}},
	}

	r := bytes.NewReader([]byte(testString))
	sc := NewScanner(r)
	res := make([]testStruct, 0)

	for sc.Scan() {
		var st testStruct

		assert.NoError(t, sc.Json(&st))
		res = append(res, st)
	}

	assert.Equal(t, expectedRes, res)
}

func TestScannerError(t *testing.T) {
	r := bytes.NewReader([]byte(testString))
	sc := NewScanner(r)
	assert.Error(t, sc.Json(&testStruct{}))
}

func TestScanner_Err(t *testing.T) {
	r := bytes.NewReader([]byte(testString))
	sc := NewScanner(r)
	assert.NoError(t, sc.Err())

	for sc.Scan() {
	}
	assert.NoError(t, sc.Err())
}
