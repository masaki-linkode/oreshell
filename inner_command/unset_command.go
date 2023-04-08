package inner_command

import (
	"fmt"
	"oreshell/ast"
	"oreshell/myvariables"
)

var CommandNameUnset = "unset"

// unsetコマンド
func Unset(simpleCommand *ast.SimpleCommand) (err error) {
	args := simpleCommand.Args()
	l := len(args)
	if l == 0 {
		return fmt.Errorf("%s: not enough arguments", CommandNameUnset)
	} else if l > 1 {
		m := map[string]string{}
		for _, v := range args {
			m[v] = ""
		}
		myvariables.Variables().AssignValuesToShellVariables(m)
	}

	return nil
}
