package ast

type Direction int

const (
	IN Direction = iota
	OUT
)

type Redirection struct {
	Direction Direction
	FdNum     int
	FilePath  string
}

type CommandSuffix struct {
	Args         []string
	Redirections []Redirection
}

type SimpleCommand struct {
	CommandName   string
	CommandSuffix CommandSuffix
}

func (me *SimpleCommand) Argv() (argv []string) {
	return append([]string{me.CommandName}, me.CommandSuffix.Args...)
}
