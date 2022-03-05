package lexer

import (
	"container/list"
	"fmt"
	"oreshell/log"
	"unicode"
	"unicode/utf8"
)

type Item struct {
	Type itemType
	Pos  int
	Val  string
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

type itemType int

const (
	ItemError itemType = iota
	ItemString
	ItemEscapeChar
	ItemQuotedString
	ItemWhitespace
	ItemEOF
	ItemRedirectionInChar
	ItemRedirectionOutChar
	ItemRedirectionFDNumChar
	ItemPipeChar
)

const eof = -1

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input   string
	state   stateFn
	pos     int
	start   int
	width   int
	lastPos int
	items   *list.List
}

func (me *Lexer) PeekItem() Item {
	log.Logger.Printf("PeekItem\n")
	item := me.items.Front().Value.(Item)
	return item
}

func (me *Lexer) NextItem() Item {
	log.Logger.Printf("NextItem\n")
	item := me.items.Front().Value.(Item)
	me.items.Remove(me.items.Front())
	me.lastPos = item.Pos
	return item
}

func Lex(input string) *Lexer {
	l := &Lexer{
		input: input,
		items: list.New(),
	}
	l.run()
	return l
}

func (me *Lexer) run() {
	log.Logger.Printf("run start\n")
	for me.state = lexText; me.state != nil; {
		me.state = me.state(me)
	}
	log.Logger.Printf("run end\n")
}

func (me *Lexer) next() rune {
	//log.Logger.Printf("next\n")

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
	//log.Logger.Printf("next r:%v\n", r)
	return r
}

func (me *Lexer) peek() rune {
	r := me.next()
	me.backup()
	return r
}

func (me *Lexer) peekpeek() rune {
	me.next()
	r := me.next()
	me.backup()
	me.backup()
	return r
}

func (me *Lexer) backup() {
	//log.Logger.Printf("backup\n")
	me.pos -= me.width
	//log.Logger.Printf("backup me.pos %d\n", me.pos)
}

func (me *Lexer) emit(t itemType) {
	log.Logger.Printf("emit\n")
	log.Logger.Printf("emit item.val:[%s]\n", me.input[me.start:me.pos])
	me.items.PushBack(Item{t, me.start, me.input[me.start:me.pos]})
	me.start = me.pos
	log.Logger.Printf("emit me.start %d\n", me.start)
}

func (me *Lexer) errorf(format string, args ...interface{}) stateFn {
	me.items.PushBack(Item{ItemError, me.start, fmt.Sprintf(format, args...)})
	return nil
}

func lexText(me *Lexer) stateFn {
	log.Logger.Printf("lexText\n")
	r := me.peek()

	if unicode.IsSpace(r) {
		return lexWhitespace
	} else if r == eof {
		me.next()
	} else if r == '\\' {
		return lexEscapeChar
	} else if r == '"' || r == '\'' {
		return lexQuotedString
	} else if r == '>' {
		return lexRedirectionOutChar
	} else if r == '<' {
		return lexRedirectionInChar
	} else if r == '0' || r == '1' || r == '2' || r == '3' || r == '4' || r == '5' || r == '6' || r == '7' || r == '8' || r == '9' {
		r2 := me.peekpeek()
		if r2 == '>' || r2 == '<' {
			return lexRedirectionFDNumChar
		} else {
			return lexString
		}
	} else if r == '|' {
		return lexPipeChar
	} else {
		return lexString
	}

	me.emit(ItemEOF)
	return nil
}

func lexWhitespace(me *Lexer) stateFn {
	log.Logger.Printf("lexWhitespace\n")
	me.next()
	for unicode.IsSpace(me.peek()) {
		me.next()
	}
	me.emit(ItemWhitespace)
	return lexText
}

func lexString(me *Lexer) stateFn {
	log.Logger.Printf("lexString\n")
	r := me.peek()
	if isDelimiter(r) {
		me.emit(ItemString)
		return lexText
	}
	me.next()
	return lexString
}

func lexEscapeChar(me *Lexer) stateFn {
	log.Logger.Printf("lexEscapeChar\n")
	me.next() // '\\'
	r := me.next()
	if r == eof {
		me.backup()
		return me.errorf("quotechar failed.")
	}
	me.emit(ItemEscapeChar)
	return lexText
}

func lexQuotedString(me *Lexer) stateFn {
	log.Logger.Printf("lexQuoteString\n")
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
	return unicode.IsSpace(r) || r == eof || r == '\\' || r == '"' || r == '\'' || r == '<' || r == '>' || r == '|'
}

func lexRedirectionInChar(me *Lexer) stateFn {
	log.Logger.Printf("lexRedirectionInChar\n")
	me.next()
	me.emit(ItemRedirectionInChar)
	return lexText
}

func lexRedirectionOutChar(me *Lexer) stateFn {
	log.Logger.Printf("lexRedirectionOutChar\n")
	me.next()
	me.emit(ItemRedirectionOutChar)
	return lexText
}

func lexRedirectionFDNumChar(me *Lexer) stateFn {
	log.Logger.Printf("lexRedirectionFDNumberChar\n")
	me.next()
	me.emit(ItemRedirectionFDNumChar)
	return lexText
}

func lexPipeChar(me *Lexer) stateFn {
	log.Logger.Printf("lexPipeChar\n")
	me.next()
	me.emit(ItemPipeChar)
	return lexText
}
