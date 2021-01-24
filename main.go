package main

import (
  "fmt"
  "bufio"
  "log"
  "os"
  "strings"
  "path/filepath"
  "io"
  "oreshell/lexer"
)

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
  //log.Printf("target %s\n", target)
  //log.Printf("targetFileName %s\n", targetFileName)

  // 指定された文字列がパスである場合
  if target != targetFileName {

    // 絶対パスの場合
    if filepath.IsAbs(target) {
      targetAbsPath = target
    // 相対パスの場合
    } else {
      targetAbsPath, err = filepath.Abs(target)
      if err != nil {
        log.Fatalf("filepath.Abs %v", err)
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
      //log.Printf("find in PATH %s\n", targetAbsPath)
      return targetAbsPath, nil
    }
  }
  return "", fmt.Errorf("%s: no such file or directory", targetFileName)
}

// cdコマンド
func chDir(words []string) (err error) {
  var dir string
  l := len(words) 
  if l == 1 {
    dir, err = os.UserHomeDir()
    if err != nil {
      log.Fatalf("os.UserHomeDir %v", err)
    }
  } else if l == 2 {
    dir = words[1]
  } else {
    return fmt.Errorf("%s: too many arguments", "cd")
  }
  return os.Chdir(dir)
}

// exitコマンド
func exit(words []string) (err error) {
  os.Exit(0)
  return nil
}

// 外部コマンドを実行する
func execExternalCommand(words []string) (err error) {
  command, err := absPathWithPATH(string(words[0]))
  if err != nil {
      fmt.Fprintln(os.Stderr, err)
      return
  }
  //log.Printf("command %s\n", command)

  // これから起動するプログラムの出力と自分の出力をつなげる
  var procAttr os.ProcAttr
  procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}

  // 該当するプログラムを探して起動する
  process, err := os.StartProcess(command, words, &procAttr)
  if err != nil {
    log.Fatalf("os.StartProcess %v", err)
  }

  // 起動したプログラムが終了するまで待つ
  _, err = process.Wait()
  if err != nil {
    log.Fatalf("process.Wait %v", err)
  }

  return nil
}

func lineToWords(line string) (words []string) {

  l := lexer.Lex(strings.Trim(line, " "))
  var word string

  for {
    token := l.NextItem()
    if token.Type == lexer.ItemWhitespace {
      words = append(words, word)
      word = ""
    } else if token.Type == lexer.ItemEOF || token.Type == lexer.ItemError {
      words = append(words, word)
      break
    } else {
      word = word + token.Unescape()
    }
  }
  //log.Printf("words: %v\n", words)

  return words
}

func main() {

  // 標準入力から文字列を読み取る準備
  reader := bufio.NewReader(os.Stdin)

  // 内部コマンド群
  internalCommands := map[string] func([]string) error {
    "cd": chDir,
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
        log.Fatalf("reader.ReadLine %v", err)
      }
    }

    // 入力文字列を空白ごとに単語に分解する
    //words := strings.Split(strings.Trim(string(line), " "), " ")
    words := lineToWords(string(line))

    // 先頭の単語に該当するコマンドを探して実行する

    // 内部コマンドか？
    internalCommand, ok := internalCommands[words[0]]
    if ok {
      // 内部コマンドを実行
      err = internalCommand(words)
    } else {
      // 外部コマンドを実行
      err = execExternalCommand(words)
    }

    if err != nil {
      fmt.Fprintln(os.Stderr, err)
    }
  }
}