package jsonl

import (
	"bytes"
	"reflect"
	"testing"
)

var testString string = "{\"name\":\"Gilbert\",\"wins\":[[\"straight\",\"7♣\"],[\"one pair\",\"10♥\"]]}\n{\"name\":\"Alexa\",\"wins\":[[\"two pair\",\"4♠\"],[\"two pair\",\"9♠\"]]}\n{\"name\":\"May\",\"wins\":[]}\n{\"name\":\"Deloise\",\"wins\":[[\"three of a kind\",\"5♣\"]]}\n"

type testStruct struct {
	Name string     `json:"name"`
	Wins [][]string `json:"wins"`
}

func TestScanner(t *testing.T) {
	expectedRes := []testStruct{
		testStruct{Name: "Gilbert", Wins: [][]string{{"straight", "7♣"}, {"one pair", "10♥"}}},
		testStruct{Name: "Alexa", Wins: [][]string{{"two pair", "4♠"}, {"two pair", "9♠"}}},
		testStruct{Name: "May", Wins: [][]string{}},
		testStruct{Name: "Deloise", Wins: [][]string{{"three of a kind", "5♣"}}},
	}

	r := bytes.NewReader([]byte(testString))
	sc := NewScanner(r)
	res := make([]testStruct, 0)

	for sc.Scan() {
		var st testStruct

		if err := sc.Json(&st); err != nil {
			t.Error(err)
		}

		res = append(res, st)
	}

	if !reflect.DeepEqual(res, expectedRes) {
		t.Errorf("%v != %v", res, expectedRes)
	}
}
