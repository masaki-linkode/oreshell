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
	rawVariables    map[string]string
	rawWords        []string
	expandVariables map[string]string
	expandedWords   []string
	redirections    []Redirection
}

func NewSimpleCommand(variables map[string]string, words []string, rs []Redirection) (me *SimpleCommand) {
	return &SimpleCommand{
		rawVariables:    variables,
		rawWords:        words,
		expandVariables: expansion.ExpandVarableValues(variables),
		expandedWords:   expansion.ExpandWords(words),
		redirections:    rs,
	}
}

func (me *SimpleCommand) Variables() map[string]string {
	return me.expandVariables
}

func (me *SimpleCommand) Argv() []string {
	return me.expandedWords
}

func (me *SimpleCommand) Args() []string {
	if len(me.expandedWords) == 0 {
		return []string{}
	}
	return me.expandedWords[1:] // 先頭以外
}

func (me *SimpleCommand) CommandName() string {
	if len(me.expandedWords) == 0 {
		return ""
	}
	return me.expandedWords[0]
}

func (me *SimpleCommand) Redirections() *[]Redirection {
	return &me.redirections
}

func (me *SimpleCommand) IsAssignVariablesOnly() bool {
	return len(me.CommandName()) == 0 && len(me.Variables()) > 0
}

type PipelineSequence struct {
	SimpleCommands []*SimpleCommand
}
