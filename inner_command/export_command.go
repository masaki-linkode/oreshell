package inner_command

import (
	"fmt"
	"oreshell/ast"
	"oreshell/myvariables"
)

// exportコマンド
func ExportEnvironmentVariable(simpleCommand *ast.SimpleCommand) (err error) {
	return doIt(simpleCommand)
}

func doIt(simpleCommand *ast.SimpleCommand) (err error) {
	commandName := "export"
	args := simpleCommand.Args()
	l := len(args)
	if l == 0 {
		return fmt.Errorf("%s: not enough arguments", commandName) // todo bashなら環境変数を一覧出力する
	} else if l == 1 {
		ok, variable_name, value := myvariables.NewAssignVariableParser().TryParse(args[0])
		if !ok {
			return fmt.Errorf("%s: %v : not a valid value", commandName, args[0]) // todo bashなら該当シェル変数を環境変数に登録する
		}
		err = myvariables.Variables().AssignValueToEnvironmentVariable(variable_name, value)
	} else {
		return fmt.Errorf("%s: %v : too many arguments", commandName, args[0])
	}
	return err
}
