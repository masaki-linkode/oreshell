package main

import (
	"bufio"
	"fmt"
	"io"
	"oreshell/ast"
	"oreshell/inner_command"
	"oreshell/lexer"
	"oreshell/log"
	"oreshell/myvariables"
	"oreshell/parser"
	"oreshell/process"
	"os"
	"strings"
)

func init() {
	var err error
	log.Logger, err = log.NewForFile("oreshell.log")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// 外部コマンドを実行する
func execExternalCommand(pipelineSequence *ast.PipelineSequence) (err error) {
	ps, err := process.NewPipelineSequence(pipelineSequence)
	if err != nil {
		return err
	}
	err = ps.Exec()
	if err != nil {
		return err
	}
	return nil
}

func main() {

	// 標準入力から文字列を読み取る準備
	reader := bufio.NewReader(os.Stdin)

	// 内部コマンド群
	internalCommands := map[string]func(*ast.SimpleCommand) error{
		inner_command.CommandNameCd:     inner_command.ChDir,
		inner_command.CommandNameExport: inner_command.ExportEnvironmentVariable,
		inner_command.CommandNameExit:   inner_command.Exit,
		inner_command.CommandNameSet:    inner_command.Set,
		inner_command.CommandNameUnset:  inner_command.Unset,
	}

	// ずっとループ
	for {
		// プロンプトを表示してユーザに入力を促す
		fmt.Printf("(ore) > ")

		// 標準入力から文字列(コマンド)を読み込む
		line, _, err := reader.ReadLine()
		if err != nil {
			// Ctrl+Dの場合
			if err == io.EOF {
				// 終了
				inner_command.Exit(nil)
			} else {
				log.Logger.Fatalf("reader.ReadLine %v", err)
			}
		}

		// 入力文字列をトリム
		trimedL := strings.Trim(string(line), " ")
		if len(trimedL) == 0 {
			continue
		}

		// 字句解析
		l := lexer.Lex(trimedL)
		// 構文解析
		pipelineSequence, err := parser.NewParser().ParsePipelineSequence(l)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		log.Logger.Printf("pipelineSequence: %+v\n", pipelineSequence)

		// 外部/内部コマンドは実行せずに、シェル変数代入のみか
		if pipelineSequence.SimpleCommands[0].IsAssignVariablesOnly() {
			myvariables.Variables().AssignValuesToShellVariables(pipelineSequence.SimpleCommands[0].Variables())
		} else {
			// 先頭の単語に該当するコマンドを探して実行する
			// 内部コマンドか？
			internalCommand, ok := internalCommands[pipelineSequence.SimpleCommands[0].CommandName()]
			if ok {
				// 内部コマンドを実行
				err = internalCommand(pipelineSequence.SimpleCommands[0])
			} else {
				// 外部コマンドを実行
				err = execExternalCommand(pipelineSequence)
			}
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
