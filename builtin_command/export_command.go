package builtin_command

import (
	"fmt"
	"oreshell/ast"
	"oreshell/myvariables"
)

var CommandNameExport = "export"

// exportコマンド
func ExportEnvironmentVariable(simpleCommand *ast.SimpleCommand) (err error) {
	return doIt(simpleCommand)
}

func doIt(simpleCommand *ast.SimpleCommand) (err error) {
	args := simpleCommand.Args()
	l := len(args)
	if l == 0 {
		return fmt.Errorf("%s: not enough arguments", CommandNameExport) // todo bashなら環境変数を一覧出力する
	} else if l == 1 {

		variable_name := ""
		value, ok := myvariables.Variables().GetValueFromShellVariables(args[0])
		if ok {
			variable_name = args[0]
		} else {
			ok, variable_name, value = myvariables.NewAssignVariableParser().TryParse(args[0])
			if !ok {
				return fmt.Errorf("%s: %v : not a valid value", CommandNameExport, args[0])
			}
		}
		err = myvariables.Variables().AssignValueToEnvironmentVariable(variable_name, value)
	} else {
		return fmt.Errorf("%s: %v : too many arguments", CommandNameExport, args[0])
	}
	return err
}
