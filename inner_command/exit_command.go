package inner_command

import (
	"oreshell/ast"
	"os"
)

var CommandNameExit = "exit"

// exitコマンド
func Exit(simpleCommand *ast.SimpleCommand) (err error) {
	os.Exit(0)
	return nil
}
