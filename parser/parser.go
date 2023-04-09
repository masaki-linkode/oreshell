package parser

import (
	"fmt"
	"oreshell/ast"
	"oreshell/lexer"
	"oreshell/log"
	"oreshell/myvariables"
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

type Parser struct {
	assignVariableParser myvariables.AssignVariableParser
}

func NewParser() Parser {
	return Parser{
		assignVariableParser: myvariables.NewAssignVariableParser(),
	}
}

func (me Parser) parseWord(l lexer.ILexer) (word string, found foundItemType, err error) {
	log.Logger.Printf("parseWord\n")
	word = ""
	for {
		token := l.PeekItem()

		if token.Type == lexer.ItemString || token.Type == lexer.ItemEscapeChar || token.Type == lexer.ItemQuotedString {
			word = word + l.NextItem().Val
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
			return "", FoundError, fmt.Errorf("ありえない %+v, %+v, %+v", token.Pos, token.Type, token.Val)
		}
	}
}

func (me Parser) parseCommandName(l lexer.ILexer) (name string, found foundItemType, err error) {
	log.Logger.Printf("parseCommandName\n")
	return me.parseWord(l)
}

func (me Parser) parseRedirection(l lexer.ILexer, fdNum int) (r *ast.Redirection, found foundItemType, err error) {
	log.Logger.Printf("parseRedirection\n")

	token := l.NextItem()
	if token.Type == lexer.ItemRedirectionInChar {
		return me.parseRedirectionIn(l, fdNum)
	} else if token.Type == lexer.ItemRedirectionOutChar {
		return me.parseRedirectionOut(l, fdNum)
	} else {
		return nil, FoundError, fmt.Errorf("ありえない")
	}
}

func (me Parser) parseRedirectionIn(l lexer.ILexer, fdNum int) (r *ast.Redirection, found foundItemType, err error) {
	log.Logger.Printf("parseRedirectionIn\n")

	for l.PeekItem().Type == lexer.ItemWhitespace {
		l.NextItem()
	}

	word, found, err := me.parseWord(l)
	if err != nil {
		return nil, FoundError, err
	}
	return ast.NewRedirection(ast.IN, fdNum, word), found, nil
}

func (me Parser) parseRedirectionOut(l lexer.ILexer, fdNum int) (r *ast.Redirection, found foundItemType, err error) {
	log.Logger.Printf("parseRedirectionOut\n")

	for l.PeekItem().Type == lexer.ItemWhitespace {
		l.NextItem()
	}

	word, found, err := me.parseWord(l)
	log.Logger.Printf("parseRedirectionOut word:%s\n", word)
	if err != nil {
		return nil, FoundError, err
	}
	return ast.NewRedirection(ast.OUT, fdNum, word), found, nil
}

func (me Parser) parseWordsAndRedirections(l lexer.ILexer) (words []string, rs []ast.Redirection, found foundItemType, err error) {
	log.Logger.Printf("parseWordsAndRedirections\n")
	r := &ast.Redirection{}

	for {
		token := l.PeekItem()

		if token.Type == lexer.ItemRedirectionFDNumChar {
			l.NextItem()
			fdNum, err := strconv.Atoi(token.Val)
			if err != nil {
				return nil, nil, FoundError, err
			}
			r, found, err = me.parseRedirection(l, fdNum)
			if err != nil {
				return nil, nil, FoundError, err
			}
			rs = append(rs, *r)
		} else if token.Type == lexer.ItemRedirectionInChar {
			l.NextItem()
			r, found, err = me.parseRedirectionIn(l, process.FD_DEFAULT_IN)
			if err != nil {
				return nil, nil, FoundError, err
			}
			rs = append(rs, *r)
		} else if token.Type == lexer.ItemRedirectionOutChar {
			l.NextItem()
			r, found, err = me.parseRedirectionOut(l, process.FD_DEFAULT_OUT)
			if err != nil {
				return nil, nil, FoundError, err
			}
			rs = append(rs, *r)
		} else {
			var word string
			word, found, err = me.parseWord(l)
			if err != nil {
				return nil, nil, FoundError, err
			}
			if len(word) > 0 {
				words = append(words, word)
			}
		}

		if found == FoundEOF || found == FoundPipe {
			log.Logger.Printf("parseWordsAndRedirections words:<%s>\n", words)
			log.Logger.Printf("parseWordsAndRedirections end\n")
			return words, rs, found, nil
		}
	}
}

func (me Parser) parseSimpleCommand(l lexer.ILexer) (s *ast.SimpleCommand, found foundItemType, err error) {
	log.Logger.Printf("parseSimpleCommand\n")

	wordsIncludeAssignVariables, rs, found, err := me.parseWordsAndRedirections(l)
	if err != nil {
		return nil, FoundError, err
	}

	variables, words := me.devideWords(wordsIncludeAssignVariables)

	log.Logger.Printf("parseSimpleCommand end\n")
	return ast.NewSimpleCommand(variables, words, rs), found, nil
}

// シェル変数の代入文群と、コマンド(とその引数)を分離する
func (me Parser) devideWords(wordsIncludeAssignVariables []string) (map[string]string, []string) {
	var words []string
	variables := map[string]string{}
	// [assign variable1] [assign variable2] ... [command]
	for i, v := range wordsIncludeAssignVariables {
		ok, variable_name, value := me.assignVariableParser.TryParse(v)
		if ok {
			log.Logger.Printf("devideWords +%v, +%v", variable_name, value)
			variables[variable_name] = value
		} else { // コマンド文字列が見つかった
			words = wordsIncludeAssignVariables[i:]
			break
		}
	}
	log.Logger.Printf("devideWords variables:+%v words:+%v\n", variables, words)
	return variables, words
}

func (me Parser) ParsePipelineSequence(l lexer.ILexer) (ps *ast.PipelineSequence, err error) {
	ps = &ast.PipelineSequence{}
	sc := &ast.SimpleCommand{}
	sc, found, err := me.parseSimpleCommand(l)
	if err != nil {
		return nil, err
	}
	ps.SimpleCommands = append(ps.SimpleCommands, sc)

	for found == FoundPipe {
		sc, found, err = me.parseSimpleCommand(l)
		if err != nil {
			return nil, err
		}
		ps.SimpleCommands = append(ps.SimpleCommands, sc)
	}

	return ps, nil
}
