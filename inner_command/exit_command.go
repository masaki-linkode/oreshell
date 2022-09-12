package inner_command

import (
	"oreshell/ast"
	"os"
)

// exitコマンド
func Exit(simpleCommand *ast.SimpleCommand) (err error) {
	os.Exit(0)
	return nil
}
