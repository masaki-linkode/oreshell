package ast

import (
	"oreshell/expansion"
)

type Direction int

const (
	IN Direction = iota
	OUT
)

type Redirection struct {
	direction Direction
	fdNum     int
	filePath  string
}

func NewRedirection(direction Direction, fdNum int, filePath string) (me *Redirection) {
	return &Redirection{direction: direction, fdNum: fdNum, filePath: filePath}
}

func (me *Redirection) Direction() Direction {
	return me.direction
}

func (me *Redirection) FdNum() int {
	return me.fdNum
}

func (me *Redirection) FilePath() string {
	return me.filePath
}

type SimpleCommand struct {
	words         []string
	expandedWords []string
	redirections  []Redirection
}

func NewSimpleCommand(words []string, rs []Redirection) (me *SimpleCommand) {
	return &SimpleCommand{words: words, expandedWords: expansion.ExpandFilenames(words), redirections: rs}
}

func (me *SimpleCommand) Argv() []string {
	return me.expandedWords
}

func (me *SimpleCommand) Args() []string {
	return me.expandedWords[1:] // 先頭以外
}

func (me *SimpleCommand) CommandName() string {
	return me.expandedWords[0]
}

func (me *SimpleCommand) Redirections() *[]Redirection {
	return &me.redirections
}

type PipelineSequence struct {
	SimpleCommands []*SimpleCommand
}
