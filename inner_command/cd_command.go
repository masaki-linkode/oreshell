package inner_command

import (
	"fmt"
	"oreshell/ast"
	"oreshell/log"
	"os"
)

var CommandNameCd = "cd"

// cdコマンド
func ChDir(simpleCommand *ast.SimpleCommand) (err error) {
	var dir string
	args := simpleCommand.Args()
	l := len(args)
	if l == 0 {
		dir, err = os.UserHomeDir()
		if err != nil {
			log.Logger.Fatalf("os.UserHomeDir %v", err)
		}
	} else if l == 1 {
		dir = args[0]
	} else {
		return fmt.Errorf("%s: too many arguments", CommandNameCd)
	}
	return os.Chdir(dir)
}
