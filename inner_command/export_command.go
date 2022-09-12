package inner_command

import (
	"fmt"
	"oreshell/ast"
	"oreshell/constdef"
	"oreshell/infra"
	"oreshell/log"
	"regexp"
	"strings"
)

// exportコマンド
func ExportEnvironmentVariable(simpleCommand *ast.SimpleCommand) (err error) {
	return newExportCommand().doIt(simpleCommand)
}

type exportCommander struct {
	osService infra.OSService
	myRegexp  *regexp.Regexp
}

func newExportCommand() exportCommander {
	return exportCommander{osService: infra.MyOSService{}, myRegexp: regexp.MustCompile(constdef.REGEX_VARIABLE_NAME)}
}

func (me exportCommander) doIt(simpleCommand *ast.SimpleCommand) (err error) {
	commandName := "export"
	args := simpleCommand.Args()
	l := len(args)
	if l == 0 {
		return fmt.Errorf("%s: not enough arguments", commandName) // todo bashなら環境変数を一覧出力する
	} else if l == 1 {
		pair := strings.SplitN(args[0], "=", 2)
		if len(pair) == 1 {
			return fmt.Errorf("%s: %v : not a valid value", commandName, args[0]) // todo bashなら該当シェル変数を環境変数に登録する
		}
		// 環境変数名が空
		if len(pair[0]) == 0 {
			return fmt.Errorf("%s: %v : not a valid identifier", commandName, args[0])
		}

		// 環境変数の文字種のチェック
		log.Logger.Printf("exportEnvironmentVariable 0: %v\n", pair[0])
		if !me.myRegexp.MatchString(pair[0]) {
			return fmt.Errorf("%s: %v : not a valid identifier", commandName, args[0])
		}

		err = me.osService.Setenv(pair[0], pair[1])
	} else {
		return fmt.Errorf("%s: %v : too many arguments", commandName, args[0])
	}
	return err
}
