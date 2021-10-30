package main

import (
	"bufio"
	"fmt"
	"io"
	"oreshell/ast"
	"oreshell/lexer"
	"oreshell/log"
	"oreshell/parser"
	"oreshell/process"
	"os"
	"path/filepath"
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

// 該当パスが存在するかどうか
func fileIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 指定された文字列が相対パスである場合、絶対パスを取得する。取得したパスが存在しなければエラーを返す。
// 指定された文字列がファイル名であるなら、環境変数PATHと連結して絶対パスを取得し存在すればそれを返す。存在しなければエラーを返す。
func absPathWithPATH(target string) (targetAbsPath string, err error) {

	// パスとファイル名を分離
	targetFileName := filepath.Base(target)
	//log.Logger.Printf("target %s\n", target)
	//log.Logger.Printf("targetFileName %s\n", targetFileName)

	// 指定された文字列がパスである場合
	if target != targetFileName {

		// 絶対パスの場合
		if filepath.IsAbs(target) {
			targetAbsPath = target
			// 相対パスの場合
		} else {
			targetAbsPath, err = filepath.Abs(target)
			if err != nil {
				log.Logger.Fatalf("filepath.Abs %v", err)
			}
		}

		if fileIsExist(targetAbsPath) {
			return targetAbsPath, nil
		} else {
			return "", fmt.Errorf("%s: no such file or directory", targetAbsPath)
		}
	}

	// 指定された文字列がファイル名である場合

	// 指定されたファイル名を環境変数パスの中から探す
	for _, path := range filepath.SplitList(os.Getenv("PATH")) {
		//log.Printf("%s\n", path)
		targetAbsPath = filepath.Join(path, targetFileName)
		if fileIsExist(targetAbsPath) {
			//log.Logger.Printf("find in PATH %s\n", targetAbsPath)
			return targetAbsPath, nil
		}
	}
	return "", fmt.Errorf("%s: no such file or directory", targetFileName)
}

// cdコマンド
func chDir(simpleCommand *ast.SimpleCommand) (err error) {
	var dir string
	l := len(simpleCommand.CommandSuffix.Args)
	if l == 0 {
		dir, err = os.UserHomeDir()
		if err != nil {
			log.Logger.Fatalf("os.UserHomeDir %v", err)
		}
	} else if l == 1 {
		dir = simpleCommand.CommandSuffix.Args[0]
	} else {
		return fmt.Errorf("%s: too many arguments", "cd")
	}
	return os.Chdir(dir)
}

// exitコマンド
func exit(simpleCommand *ast.SimpleCommand) (err error) {
	os.Exit(0)
	return nil
}

// 外部コマンドを実行する
func execExternalCommand(simpleCommand *ast.SimpleCommand) (err error) {
	command, err := absPathWithPATH(string(simpleCommand.CommandName))
	if err != nil {
		return err
	}
	log.Logger.Printf("command %s\n", command)
	log.Logger.Printf("args : %v", simpleCommand.CommandSuffix.Args)

	var procAttr os.ProcAttr
	procAttr.Files, err = process.CreateProcAttrFiles(&simpleCommand.CommandSuffix.Redirections)
	if err != nil {
		return err
	}

	// 該当するプログラムを探して起動する
	process, err := os.StartProcess(command, simpleCommand.Argv(), &procAttr)
	if err != nil {
		log.Logger.Fatalf("os.StartProcess %v", err)
	}

	// 起動したプログラムが終了するまで待つ
	_, err = process.Wait()
	if err != nil {
		log.Logger.Fatalf("process.Wait %v", err)
	}

	return nil
}

func main() {

	// 標準入力から文字列を読み取る準備
	reader := bufio.NewReader(os.Stdin)

	// 内部コマンド群
	internalCommands := map[string]func(*ast.SimpleCommand) error{
		"cd":   chDir,
		"exit": exit,
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
				exit(nil)
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
		simpleCommand, err := parser.ParseSimpleCommand(l)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		log.Logger.Printf("simpleCommand: %v\n", simpleCommand)

		// 先頭の単語に該当するコマンドを探して実行する
		// 内部コマンドか？
		internalCommand, ok := internalCommands[simpleCommand.CommandName]
		if ok {
			// 内部コマンドを実行
			err = internalCommand(simpleCommand)
		} else {
			// 外部コマンドを実行
			err = execExternalCommand(simpleCommand)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
