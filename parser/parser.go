package parser

import (
	"fmt"
	"oreshell/ast"
	"oreshell/lexer"
	"oreshell/log"
	"oreshell/process"
	"strconv"
)

func parseWord(l *lexer.Lexer) (word string, foundEOF bool, err error) {
	log.Logger.Printf("parseWord\n")
	word = ""
	for {
		token := l.PeekItem()

		if token.Type == lexer.ItemString || token.Type == lexer.ItemEscapeChar || token.Type == lexer.ItemQuotedString {
			word = word + l.NextItem().Unescape()
			log.Logger.Printf("word: %s\n", word)
		} else if token.Type == lexer.ItemWhitespace {
			log.Logger.Printf("found whitespace %s\n", word)
			l.NextItem()
			return word, false, nil
		} else if token.Type == lexer.ItemRedirectionInChar || token.Type == lexer.ItemRedirectionOutChar {
			log.Logger.Printf("found redirectionchar\n")
			return word, false, nil
		} else if token.Type == lexer.ItemEOF {
			log.Logger.Printf("found eof %s\n", word)
			l.NextItem()
			if len(word) == 0 {
				return "", true, fmt.Errorf("構文エラー")
			}
			return word, true, nil
		} else {
			return "", false, fmt.Errorf("ありえない")
		}
	}
}

func parseCommandName(l *lexer.Lexer) (name string, foundEOF bool, err error) {
	log.Logger.Printf("parseCommandName\n")
	return parseWord(l)
}

func parseRedirection(l *lexer.Lexer, fdNum int) (r *ast.Redirection, foundEOF bool, err error) {
	log.Logger.Printf("parseRedirection\n")
	r = &ast.Redirection{}
	token := l.NextItem()
	if token.Type == lexer.ItemRedirectionInChar {
		return parseRedirectionIn(l, fdNum)
	} else if token.Type == lexer.ItemRedirectionOutChar {
		return parseRedirectionOut(l, fdNum)
	} else {
		return nil, false, fmt.Errorf("ありえない")
	}
}

func parseRedirectionIn(l *lexer.Lexer, fdNum int) (r *ast.Redirection, foundEOF bool, err error) {
	log.Logger.Printf("parseRedirectionIn\n")
	r = &ast.Redirection{}

	for l.PeekItem().Type == lexer.ItemWhitespace {
		l.NextItem()
	}

	word, foundEOF, err := parseWord(l)
	if err != nil {
		return nil, false, err
	}
	return &ast.Redirection{Direction: ast.IN, FdNum: fdNum, FilePath: word}, foundEOF, nil
}

func parseRedirectionOut(l *lexer.Lexer, fdNum int) (r *ast.Redirection, foundEOF bool, err error) {
	log.Logger.Printf("parseRedirectionOut\n")
	r = &ast.Redirection{}

	for l.PeekItem().Type == lexer.ItemWhitespace {
		l.NextItem()
	}

	word, foundEOF, err := parseWord(l)
	log.Logger.Printf("parseRedirectionOut word:%s\n", word)
	if err != nil {
		return nil, false, err
	}
	return &ast.Redirection{Direction: ast.OUT, FdNum: fdNum, FilePath: word}, foundEOF, nil
}

func parseCommandSuffix(l *lexer.Lexer) (c *ast.CommandSuffix, foundEOF bool, err error) {
	log.Logger.Printf("parseCommandSuffix\n")
	r := &ast.Redirection{}
	var args []string
	var rs []ast.Redirection

	for {
		token := l.PeekItem()

		if token.Type == lexer.ItemRedirectionFDNumChar {
			l.NextItem()
			fdNum, err := strconv.Atoi(token.Val)
			if err != nil {
				return nil, false, err
			}
			r, foundEOF, err = parseRedirection(l, fdNum)
			if err != nil {
				return nil, false, err
			}
			rs = append(rs, *r)
		} else if token.Type == lexer.ItemRedirectionInChar {
			l.NextItem()
			r, foundEOF, err = parseRedirectionIn(l, process.FD_DEFAULT_IN)
			if err != nil {
				return nil, false, err
			}
			rs = append(rs, *r)
		} else if token.Type == lexer.ItemRedirectionOutChar {
			l.NextItem()
			r, foundEOF, err = parseRedirectionOut(l, process.FD_DEFAULT_OUT)
			if err != nil {
				return nil, false, err
			}
			rs = append(rs, *r)
		} else {
			var word string
			word, foundEOF, err = parseWord(l)
			if err != nil {
				return nil, false, err
			}
			args = append(args, word)
		}

		if foundEOF {
			c = &ast.CommandSuffix{Args: args, Redirections: rs}
			return c, true, nil
		}
	}
}

func ParseSimpleCommand(l *lexer.Lexer) (s *ast.SimpleCommand, err error) {

	s = &ast.SimpleCommand{}

	word, foundEOF, err := parseWord(l)
	if err != nil {
		return nil, err
	}
	s.CommandName = word

	if !foundEOF {
		c, _, err := parseCommandSuffix(l)
		if err != nil {
			return nil, err
		}
		s.CommandSuffix = *c
	}

	return s, nil
}
