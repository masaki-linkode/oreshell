package myvariables

import (
	"fmt"
	"oreshell/log"
	"testing"
)

func init() {
	log.Logger = log.New()
}

type testItem struct {
	name   string
	input  string
	expect parseResult
}

type parseResult struct {
	ok            bool
	variable_name string
	value         string
}

var testItems = []testItem{
	{"01.OK", "HOGE=HIGE", parseResult{true, "HOGE", "HIGE"}},
	{"02.OK 値がない", "HOGE=", parseResult{true, "HOGE", ""}},
	{"03.OK 値がクォートされている", "HOGE=\"HIGE\"", parseResult{true, "HOGE", "\"HIGE\""}},
	{"11.ERR 変数名がない", "=HIGE", parseResult{false, "", ""}},
	{"12.ERR 「=」がクォートされている", "HOGE\"=HIGE\"", parseResult{false, "", ""}},
	{"13.ERR 変数名がクォートされている", "\"HOGE\"=HIGE", parseResult{false, "", ""}},
	{"13.ERR 変数名の一部がクォートされている", "HO\"GE\"=HIGE", parseResult{false, "", ""}},
}

func parse(t *testItem) parseResult {
	ok, variable_name, value := NewAssignVariableParser().TryParse(t.input)
	return parseResult{ok, variable_name, value}
}

func equal(i1, i2 parseResult) bool {
	if i1.ok == i2.ok && i1.variable_name == i2.variable_name && i1.value == i2.value {
		return true
	}
	return false
}

func Test(t *testing.T) {
	for _, test := range testItems {
		fmt.Printf("%s\n", test.name)
		actual := parse(&test)
		if !equal(test.expect, actual) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%+v", test.name, actual, test.expect)
		}
	}
}
