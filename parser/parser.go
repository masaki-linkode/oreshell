package parser

import (
	"fmt"
	"oreshell/ast"
	"oreshell/lexer"
	"oreshell/log"
	"oreshell/process"
	"strconv"
)

type foundItemType int

const (
	FoundEOF foundItemType = iota
	FoundPipe
	FoundError
	Other
)

func parseWord(l *lexer.Lexer) (word string, found foundItemType, err error) {
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
			for l.PeekItem().Type == lexer.ItemWhitespace {
				l.NextItem()
			}
			return word, Other, nil
		} else if token.Type == lexer.ItemRedirectionInChar || token.Type == lexer.ItemRedirectionOutChar {
			log.Logger.Printf("found redirectionchar\n")
			return word, Other, nil
		} else if token.Type == lexer.ItemPipeChar {
			log.Logger.Printf("found pipechar\n")
			l.NextItem()
			for l.PeekItem().Type == lexer.ItemWhitespace {
				l.NextItem()
			}
			return word, FoundPipe, nil
		} else if token.Type == lexer.ItemEOF {
			log.Logger.Printf("found eof %s\n", word)
			l.NextItem()
			return word, FoundEOF, nil
		} else {
			return "", FoundError, fmt.Errorf("ありえない")
		}
	}
}

func parseCommandName(l *lexer.Lexer) (name string, found foundItemType, err error) {
	log.Logger.Printf("parseCommandName\n")
	return parseWord(l)
}

func parseRedirection(l *lexer.Lexer, fdNum int) (r *ast.Redirection, found foundItemType, err error) {
	log.Logger.Printf("parseRedirection\n")
	r = &ast.Redirection{}
	token := l.NextItem()
	if token.Type == lexer.ItemRedirectionInChar {
		return parseRedirectionIn(l, fdNum)
	} else if token.Type == lexer.ItemRedirectionOutChar {
		return parseRedirectionOut(l, fdNum)
	} else {
		return nil, FoundError, fmt.Errorf("ありえない")
	}
}

func parseRedirectionIn(l *lexer.Lexer, fdNum int) (r *ast.Redirection, found foundItemType, err error) {
	log.Logger.Printf("parseRedirectionIn\n")
	r = &ast.Redirection{}

	for l.PeekItem().Type == lexer.ItemWhitespace {
		l.NextItem()
	}

	word, found, err := parseWord(l)
	if err != nil {
		return nil, FoundError, err
	}
	return &ast.Redirection{Direction: ast.IN, FdNum: fdNum, FilePath: word}, found, nil
}

func parseRedirectionOut(l *lexer.Lexer, fdNum int) (r *ast.Redirection, found foundItemType, err error) {
	log.Logger.Printf("parseRedirectionOut\n")
	r = &ast.Redirection{}

	for l.PeekItem().Type == lexer.ItemWhitespace {
		l.NextItem()
	}

	word, found, err := parseWord(l)
	log.Logger.Printf("parseRedirectionOut word:%s\n", word)
	if err != nil {
		return nil, FoundError, err
	}
	return &ast.Redirection{Direction: ast.OUT, FdNum: fdNum, FilePath: word}, found, nil
}

func parseCommandSuffix(l *lexer.Lexer) (c *ast.CommandSuffix, found foundItemType, err error) {
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
				return nil, FoundError, err
			}
			r, found, err = parseRedirection(l, fdNum)
			if err != nil {
				return nil, FoundError, err
			}
			rs = append(rs, *r)
		} else if token.Type == lexer.ItemRedirectionInChar {
			l.NextItem()
			r, found, err = parseRedirectionIn(l, process.FD_DEFAULT_IN)
			if err != nil {
				return nil, FoundError, err
			}
			rs = append(rs, *r)
		} else if token.Type == lexer.ItemRedirectionOutChar {
			l.NextItem()
			r, found, err = parseRedirectionOut(l, process.FD_DEFAULT_OUT)
			if err != nil {
				return nil, FoundError, err
			}
			rs = append(rs, *r)
		} else {
			var word string
			word, found, err = parseWord(l)
			if err != nil {
				return nil, FoundError, err
			}
			if len(word) > 0 {
				args = append(args, word)
			}
		}

		if found == FoundEOF || found == FoundPipe {
			log.Logger.Printf("parseCommandSuffix args:<%s>\n", args)
			c = &ast.CommandSuffix{Args: args, Redirections: rs}

			log.Logger.Printf("parseCommandSuffix end\n")
			return c, found, nil
		}
	}
}

func parseSimpleCommand(l *lexer.Lexer) (s *ast.SimpleCommand, found foundItemType, err error) {
	log.Logger.Printf("parseSimpleCommand\n")

	s = &ast.SimpleCommand{}

	word, found, err := parseCommandName(l)
	if err != nil {
		return nil, FoundError, err
	}
	log.Logger.Printf("s.CommandName:<%s>\n", word)
	s.CommandName = word

	if found != FoundEOF && found != FoundPipe {
		c, f, err := parseCommandSuffix(l)
		if err != nil {
			return nil, FoundError, err
		}
		s.CommandSuffix = *c
		found = f
	}

	log.Logger.Printf("parseSimpleCommand end\n")
	return s, found, nil
}

func ParsePipelineSequence(l *lexer.Lexer) (ps *ast.PipelineSequence, err error) {
	ps = &ast.PipelineSequence{}
	sc := &ast.SimpleCommand{}
	sc, found, err := parseSimpleCommand(l)
	if err != nil {
		return nil, err
	}
	ps.SimpleCommands = append(ps.SimpleCommands, sc)

	for found == FoundPipe {
		sc, found, err = parseSimpleCommand(l)
		if err != nil {
			return nil, err
		}
		ps.SimpleCommands = append(ps.SimpleCommands, sc)
	}

	return ps, nil
}
