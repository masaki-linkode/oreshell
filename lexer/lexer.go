package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Item struct {
	Type itemType
	Pos int
	Val string
}

func (me Item) String() string {
	switch {
	case me.Type == ItemEOF:
		return "EOF"
	case me.Type == ItemError:
		return me.Val
	}
	return fmt.Sprintf("%q", me.Val)
}

func (me Item) Unescape() string {
	switch {
	case me.Type == ItemEscapeChar:
		return string(me.Val[1])
	case me.Type == ItemQuotedString:
		return strings.Trim(me.Val, string(me.Val[0]))
	}
	return me.Val
}

type itemType int

const (
	ItemError itemType = iota
	ItemString
	ItemEscapeChar
	ItemQuotedString
	ItemWhitespace
	ItemEOF
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input   string
	state   stateFn
	pos     int
	start   int
	width   int
	lastPos int
	items   chan Item
}

func (me *lexer) NextItem() Item {
	item := <-me.items
	me.lastPos = item.Pos
	return item
}

func Lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan Item),
	}
	go l.run()
	return l
}

func (me *lexer) run() {
	for me.state = lexText; me.state != nil; {
		me.state = me.state(me)
	}
}

func (me *lexer) next() rune {
	//fmt.Printf("next\n")

	// 現在位置が入力文字列全体を超えたか
	if int(me.pos) >= len(me.input) {
		me.width = 0
		return eof
	}
	// 現在位置の文字を取り出す
	r, w := utf8.DecodeRuneInString(me.input[me.pos:])
	me.width = w
	// 現在位置を文字幅分だけ進める
	me.pos += w
	//fmt.Printf("next r:%v\n", r)
	return r
}

func (me *lexer) peek() rune {
	r := me.next()
	me.backup()
	return r
}

func (me *lexer) backup() {
	//fmt.Printf("backup\n")
	me.pos -= me.width
	//fmt.Printf("backup me.pos %d\n", me.pos)
}

func (me *lexer) emit(t itemType) {
	//fmt.Printf("emit\n")
	//fmt.Printf("emit item.val:[%s]\n", me.input[me.start:me.pos])
	me.items <- Item{t, me.start, me.input[me.start:me.pos]}
	me.start = me.pos
	//fmt.Printf("emit me.start %d\n", me.start)
}

func (me *lexer) errorf(format string, args ...interface{}) stateFn {
	me.items <- Item{ItemError, me.start, fmt.Sprintf(format, args...)}
	return nil
}

func lexText(me *lexer) stateFn {
	//fmt.Printf("lexText\n")
	r := me.peek()
				
	if unicode.IsSpace(r) {
		return lexWhitespace
	} else if r == eof {
		me.next()
	} else if r == '\\' {
		return lexEscapeChar
	} else if r == '"' || r == '\'' {
		return lexQuotedString
	} else {
		return lexString
	}
	
	me.emit(ItemEOF)
	return nil
}

func lexWhitespace(me *lexer) stateFn {
	//fmt.Printf("lexWhitespace\n")
	me.next()
	for unicode.IsSpace(me.peek()) {
		me.next()
	}
	me.emit(ItemWhitespace)
	return lexText
}

func lexString(me *lexer) stateFn {
	//fmt.Printf("lexString\n")
	r := me.peek()
	if isDelimiter(r) {
		me.emit(ItemString)
		return lexText
	}
	me.next()
	return lexString
}

func lexEscapeChar(me *lexer) stateFn {
	//fmt.Printf("lexEscapeChar\n")
	me.next() // '\\'
	r := me.next()
	if r == eof {
		me.backup()
		return me.errorf("quotechar failed.")
	}
	me.emit(ItemEscapeChar)
	return lexText
}

func lexQuotedString(me *lexer) stateFn {
	//fmt.Printf("lexQuoteString\n")
	q := me.next() // '"' or '\''

	for {
		r := me.next()
		if r == eof {
			me.backup()
			return me.errorf("quotestring failed.")
		} else if r == q {
			break
		}
	}
	me.emit(ItemQuotedString)
	return lexText
}

func isDelimiter(r rune) bool {
	return unicode.IsSpace(r) || r == eof || r == '\\' || r == '"' || r == '\''
}