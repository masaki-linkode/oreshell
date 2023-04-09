package parser

import (
	"container/list"
	"fmt"
	"oreshell/lexer"
	"oreshell/log"
	"testing"
)

func init() {
	log.Logger = log.New()
}

type parseWordTest struct {
	name         string
	items        []lexer.Item
	expect_word  string
	expect_found foundItemType
	expect_err   error
}

var parseWordTests = []parseWordTest{
	{
		name: "01 最後までスペースを読み飛ばす",
		items: []lexer.Item{
			{
				Type: lexer.ItemWhitespace,
				Val:  " ",
			},
			{
				Type: lexer.ItemWhitespace,
				Val:  " ",
			},
			{
				Type: lexer.ItemEOF,
				Val:  "",
			},
		},
		expect_word:  "",
		expect_found: Other,
		expect_err:   nil,
	},
	{
		name: "02 単語が見つかるまでスペースを読み飛ばす",
		items: []lexer.Item{
			{
				Type: lexer.ItemWhitespace,
				Val:  " ",
			},
			{
				Type: lexer.ItemWhitespace,
				Val:  " ",
			},
			{
				Type: lexer.ItemString,
				Val:  " ",
			},
			{
				Type: lexer.ItemEOF,
				Val:  "",
			},
		},
		expect_word:  "",
		expect_found: Other,
		expect_err:   nil,
	},
}

type testlexer struct {
	items *list.List
}

func newTestLexer(items []lexer.Item) *testlexer {
	me := &testlexer{
		items: list.New(),
	}
	for _, item := range items {
		me.items.PushBack(item)
	}
	return me
}

func (me *testlexer) PeekItem() lexer.Item {
	log.Logger.Printf("PeekItem\n")
	item := me.items.Front().Value.(lexer.Item)
	return item
}

func (me *testlexer) NextItem() lexer.Item {
	log.Logger.Printf("NextItem\n")
	item := me.items.Front().Value.(lexer.Item)
	me.items.Remove(me.items.Front())
	return item
}

func TestParseWord(t *testing.T) {
	for _, test := range parseWordTests {
		fmt.Printf("%s\n", test.name)

		word, found, err := NewParser().parseWord(newTestLexer(test.items))

		if word != test.expect_word {
			t.Errorf("%s: word got\n\t%+v\nexpected\n\t%v", test.name, word, test.expect_word)
		} else if found != test.expect_found {
			t.Errorf("%s: found got\n\t%+v\nexpected\n\t%v", test.name, found, test.expect_found)
		} else if err != test.expect_err {
			t.Errorf("%s: err got\n\t%+v\nexpected\n\t%v", test.name, err, test.expect_err)
		}
	}
}

type expectSimpleCommand struct {
	variables    map[string]string
	command_name string
	args         []string
}

type parsePipelineSequenceTest struct {
	name                 string
	input                string
	expectSimpleCommands []expectSimpleCommand
	expect_err           error
}

var parsePipelineSequenceTests = []parsePipelineSequenceTest{
	{
		name:  ".引数0",
		input: "test",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test",
				args:         []string{},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".引数1",
		input: "test abc",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test",
				args:         []string{"abc"},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".引数2",
		input: "test abc def",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test",
				args:         []string{"abc", "def"},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".コマンドの前に空白",
		input: " test abc def",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test",
				args:         []string{"abc", "def"},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".最後の引数の後に空白",
		input: "test abc def ",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test",
				args:         []string{"abc", "def"},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".パイプ 引数なし|引数なし",
		input: "test1|test2",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test1",
				args:         []string{},
			},
			{
				command_name: "test2",
				args:         []string{},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".パイプ 引数あり|引数あり",
		input: "test1 abc def|test2 ghi jkl",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test1",
				args:         []string{"abc", "def"},
			},
			{
				command_name: "test2",
				args:         []string{"ghi", "jkl"},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".パイプの前に空白 引数あり |引数あり",
		input: "test1 abc def |test2 ghi jkl",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test1",
				args:         []string{"abc", "def"},
			},
			{
				command_name: "test2",
				args:         []string{"ghi", "jkl"},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".パイプの後に空白 引数あり| 引数あり",
		input: "test1 abc def| test2 ghi jkl",
		expectSimpleCommands: []expectSimpleCommand{
			{
				command_name: "test1",
				args:         []string{"abc", "def"},
			},
			{
				command_name: "test2",
				args:         []string{"ghi", "jkl"},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".パイプの前後に空白 引数あり | 引数あり",
		input: "test1 abc def | test2 ghi jkl",
		expectSimpleCommands: []expectSimpleCommand{{
			command_name: "test1",
			args:         []string{"abc", "def"},
		},
			{
				command_name: "test2",
				args:         []string{"ghi", "jkl"},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".シェル変数代入",
		input: "HOGE=HIGE",
		expectSimpleCommands: []expectSimpleCommand{
			{
				variables:    map[string]string{"HOGE": "HIGE"},
				command_name: "",
				args:         []string{},
			},
		},
		expect_err: nil,
	},
	{
		name:  ".シェル変数代入 複数",
		input: "HOGE=HIGE HUGE=HEGE",
		expectSimpleCommands: []expectSimpleCommand{
			{
				variables:    map[string]string{"HOGE": "HIGE", "HUGE": "HEGE"},
				command_name: "",
				args:         []string{},
			},
		},
		expect_err: nil,
	},
}

func TestParsePipelineSequence(t *testing.T) {
	for _, test := range parsePipelineSequenceTests {
		fmt.Printf("%s\n", test.name)

		ps, err := NewParser().ParsePipelineSequence(lexer.Lex(test.input))

		if err != test.expect_err {
			t.Errorf("%s: err got\n\t%+v\nexpected\n\t%v", test.name, err, test.expect_err)
		}

		if len(ps.SimpleCommands) != len(test.expectSimpleCommands) {
			t.Errorf("%s: len(SimpleCommands) got\n\t%+v\nexpected\n\t%v", test.name, len(ps.SimpleCommands), len(test.expectSimpleCommands))
		}

		for i, expect := range test.expectSimpleCommands {
			actual := ps.SimpleCommands[i]
			if actual.CommandName() != expect.command_name {
				t.Errorf("%s: %d commandName got\n\t%+v\nexpected\n\t%+v", test.name, i, actual.CommandName(), expect.command_name)
			}
			if len(actual.Args()) != len(expect.args) {
				t.Errorf("%s: %d Args() got\n\t%+v\nexpected\n\t%v", test.name, i, len(actual.Args()), len(expect.args))
			}
			for j, actualArg := range actual.Args() {
				expectArg := expect.args[j]
				if actualArg != expectArg {
					t.Errorf("%s: %d  Args[%d] got\n\t%+v\nexpected\n\t%v", test.name, i, j, actualArg, expectArg)
				}
			}
			mapEqual(actual.Variables(), expect.variables)
		}
	}
}

func mapEqual(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}

	return true
}
