package jsonUnsafe

import (
	"reflect"
	"testing"

	"github.com/ElDelto/testutil"
)

type testStruct struct {
	unexported string
	Exported   string
}

var (
	ts       = testStruct{"val0", "val1"}
	testJSON = []byte(`{
                  "unexported": "val0", 
                  "Exported":   "val1"
                }`)

	tm = map[string]interface{}{
		"unexported": "val0",
		"Exported":   "val1",
	}
)

func TestMarshal(t *testing.T) {
	json, err := Marshal(&ts)

	testutil.CheckError(t, err)
	testutil.ShouldBeEqual(t, "JSON", testJSON, json)
}

func TestUnmarshal(t *testing.T) {
	tsLocal := testStruct{}
	err := Unmarshal(testJSON, &tsLocal)

	if err != nil {
		testutil.Error(t, err)
	}

	if tsLocal != ts {
		testutil.ShouldBe(t, "struct", ts, tsLocal)
	}

	m := map[string]interface{}{}
	err = Unmarshal(testJSON, m)

	if err != nil {
		testutil.Error(t, err)
	}

	if reflect.DeepEqual(m, tm) {
		testutil.ShouldBe(t, "map", tm, m)
	}
}

func TestFindValue(t *testing.T) {
	m := map[string]interface{}{
		"value0": 0,
		"Value1": 1,
		"valUe2": 2,
	}

	testData := []struct {
		in  string
		out int
	}{
		{"value0", 0},
		{"value1", 1},
		{"value2", 2},
	}

	for _, td := range testData {
		v, _ := findValue(m, td.in)
		if v != td.out {
			testutil.ShouldBe(t, "value", td.out, v)
		}
	}
}

func TestSetUnexportedField(t *testing.T) {
	tsCopy := ts
	rtValue := reflect.ValueOf(&tsCopy).Elem()
	rtField := rtValue.Field(0)

	testString := "TEST"
	_ = setUnexportedField(&rtField, testString)
	if tsCopy.unexported != testString {
		testutil.ShouldBe(t, "unexported field", testString, tsCopy.unexported)
	}

	testInt := 10
	err := setUnexportedField(&rtField, testInt)
	if err == nil {
		testutil.ShouldBe(t, "error", "available", nil)
	}
}
