package expansion

import (
	"fmt"
	"oreshell/log"
	"reflect"
	"testing"
)

func init() {
	log.Logger = log.New()
}

type testItem struct {
	name   string
	input  string
	expect result
}

type result struct {
	files []string
}

var testItems = []testItem{
	{"01.", "*_test.go", result{[]string{"filename_exp_test.go", "shell_parameter_exp_test.go"}}},
	{"02.", "testdata/\"a.txt\"", result{[]string{"testdata/\"a.txt\""}}},
	{"03.", "testdata/a.txt", result{[]string{}}},
	{"04.", "testdata/\"*\"", result{[]string{"testdata/\"a.txt\""}}}, // bashではこのファイル名を取得しない。todo
	{"05.", "testdata/\\\"*\\\"", result{[]string{"testdata/\"a.txt\""}}},
}

func doTest(t *testItem) result {
	files := expandFilename(t.input)
	return result{files}
}

func equal(i1, i2 result) bool {
	if (len(i1.files) == 0 && len(i2.files) == 0) || reflect.DeepEqual(i1.files, i2.files) {
		return true
	}
	return false
}

func Test(t *testing.T) {
	for _, test := range testItems {
		fmt.Printf("%s\n", test.name)
		actual := doTest(&test)
		if !equal(test.expect, actual) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%+v", test.name, actual, test.expect)
		}
	}
}
