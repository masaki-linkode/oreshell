package inner_command

import (
	"fmt"
	"oreshell/ast"
	"oreshell/myvariables"
	"os"
)

var CommandNameSet = "set"

// setコマンド
func Set(simpleCommand *ast.SimpleCommand) (err error) {
	args := simpleCommand.Args()
	l := len(args)
	if (l) == 0 {
		for kv := range myvariables.GetKVIterator() {
			fmt.Fprintf(os.Stdout, "%s=%s\n", kv.VariableName, kv.Value)
		}
		return nil
	} else { // 現時点では引数を取らない
		return fmt.Errorf("%s: %v : too many arguments", simpleCommand.CommandName(), args[0])
	}
}
