package lexer

import (
	"fmt"
	"oreshell/log"
	"testing"
)

func init() {
	log.Logger = log.New()
}

type lexTest struct {
	name  string
	input string
	items []Item
}

var tEOF = Item{ItemEOF, 0, ""}

var lexTests = []lexTest{
	{"01.空", "", []Item{tEOF}},
	{"02.空白扱い文字", " \t\n", []Item{{ItemWhitespace, 0, " \t\n"}, tEOF}},
	{"03.1文字", `a`, []Item{
		{ItemString, 0, `a`},
		tEOF,
	}},
	{"04.strings", `a b`, []Item{
		{ItemString, 0, `a`},
		{ItemWhitespace, 0, ` `},
		{ItemString, 0, `b`},
		tEOF,
	}},
	{"05.escape char", `a \ b`, []Item{
		{ItemString, 0, `a`},
		{ItemWhitespace, 0, ` `},
		{ItemEscapeChar, 0, "\\ "},
		{ItemString, 0, `b`},
		tEOF,
	}},
	{"06.quoted string", `a "b c" e"f g"h`, []Item{
		{ItemString, 0, `a`},
		{ItemWhitespace, 0, ` `},
		{ItemQuotedString, 0, `"b c"`},
		{ItemWhitespace, 0, ` `},
		{ItemString, 0, `e`},
		{ItemQuotedString, 0, `"f g"`},
		{ItemString, 0, `h`},
		tEOF,
	}},
	{"07.redirection out", `a >b c>d > e`, []Item{
		{ItemString, 0, `a`},
		{ItemWhitespace, 0, ` `},
		{ItemRedirectionOutChar, 0, `>`},
		{ItemString, 0, `b`},
		{ItemWhitespace, 0, ` `},
		{ItemString, 0, `c`},
		{ItemRedirectionOutChar, 0, `>`},
		{ItemString, 0, `d`},
		{ItemWhitespace, 0, ` `},
		{ItemRedirectionOutChar, 0, `>`},
		{ItemWhitespace, 0, ` `},
		{ItemString, 0, `e`},
		tEOF,
	}},
	{"08.redirection in", `a <b c<d < e`, []Item{
		{ItemString, 0, `a`},
		{ItemWhitespace, 0, ` `},
		{ItemRedirectionInChar, 0, `<`},
		{ItemString, 0, `b`},
		{ItemWhitespace, 0, ` `},
		{ItemString, 0, `c`},
		{ItemRedirectionInChar, 0, `<`},
		{ItemString, 0, `d`},
		{ItemWhitespace, 0, ` `},
		{ItemRedirectionInChar, 0, `<`},
		{ItemWhitespace, 0, ` `},
		{ItemString, 0, `e`},
		tEOF,
	}},
	{"09.redirection fdnum", `a 1<b c2<d 3< e`, []Item{
		{ItemString, 0, `a`},
		{ItemWhitespace, 0, ` `},
		{ItemRedirectionFDNumChar, 0, `1`},
		{ItemRedirectionInChar, 0, `<`},
		{ItemString, 0, `b`},
		{ItemWhitespace, 0, ` `},
		{ItemString, 0, `c2`},
		{ItemRedirectionInChar, 0, `<`},
		{ItemString, 0, `d`},
		{ItemWhitespace, 0, ` `},
		{ItemRedirectionFDNumChar, 0, `3`},
		{ItemRedirectionInChar, 0, `<`},
		{ItemWhitespace, 0, ` `},
		{ItemString, 0, `e`},
		tEOF,
	}},
}

func collect(t *lexTest) (items []Item) {
	l := Lex(t.input)
	for {
		item := l.NextItem()
		items = append(items, item)
		if item.Type == ItemEOF || item.Type == ItemError {
			break
		}
	}
	return
}

func equal(i1, i2 []Item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].Type != i2[k].Type {
			return false
		}
		if i1[k].Val != i2[k].Val {
			return false
		}
		if checkPos && i1[k].Pos != i2[k].Pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		fmt.Printf("%s\n", test.name)
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
