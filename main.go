package main

import (
  "fmt"
  "bufio"
  "log"
  "os"
  "strings"
  "path/filepath"
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
  log.Printf("target %s\n", target)
  log.Printf("targetFileName %s\n", targetFileName)

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
    log.Printf("%s\n", path)
    targetAbsPath = filepath.Join(path, targetFileName)
    if fileIsExist(targetAbsPath) {
      log.Printf("find in PATH %s\n", targetAbsPath)
      return targetAbsPath, nil
    }
  }
  return "", fmt.Errorf("%s: no such file or directory", targetFileName)
}

func main() {

  // 標準入力から文字列を読み取る準備
  reader := bufio.NewReader(os.Stdin)

  // 0.ずっとループ
  for {
    // 1.プロンプトを表示してユーザに入力を促す
    fmt.Printf("(ore) > ")
    
    // 3.標準入力から文字列(コマンド)を読み込む
    line, _, err := reader.ReadLine()
    if err != nil {
      log.Fatalf("ReadLine %v", err)
    }
    words := strings.Split(string(line), " ")

    command, err := absPathWithPATH(string(words[0]))
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        continue
    }
    log.Printf("command %s\n", command)

    // これから起動するプログラムの出力と自分の出力をつなげる(6,7)
    var procAttr os.ProcAttr
    procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}

    // 4.入力文字列に該当するプログラムを探して起動する
    process, err := os.StartProcess(command, words, &procAttr)

    if err != nil {
      log.Fatalf("StartProcess %v", err)
    }

    // 起動したプログラムが終了するまで待つ(8を待つ)
    _, err = process.Wait()
    if err != nil {
      log.Fatalf("Wait %v", err)
    }
  }
}